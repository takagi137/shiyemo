package com.shiyemo.app.ui.profile

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.local.prefs.TokenManager
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.model.dto.UserDto
import com.shiyemo.app.data.repository.ProfileRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

/**
 * 个人中心页面 Tab
 */
enum class ProfileTab(val title: String) {
    POSTS("帖子"),
    BOOKMARKS("收藏"),
    FOLLOWING("关注"),
    FOLLOWERS("粉丝")
}

/**
 * 个人中心 UI 状态
 */
data class ProfileUiState(
    val user: UserDto? = null,
    val posts: List<PostDto> = emptyList(),
    val bookmarks: List<PostDto> = emptyList(),
    val followingList: List<UserDto> = emptyList(),
    val followersList: List<UserDto> = emptyList(),
    val isLoading: Boolean = false,
    val isLoadingPosts: Boolean = false,
    val isLoadingBookmarks: Boolean = false,
    val isLoadingFollowing: Boolean = false,
    val isLoadingFollowers: Boolean = false,
    val error: String? = null,
    val isLoggedIn: Boolean = false,
    val isOwnProfile: Boolean = true, // 是否是自己的主页
    val isFollowing: Boolean = false, // 是否已关注该用户
    val selectedTab: ProfileTab = ProfileTab.POSTS,
    val postsPage: Int = 1,
    val bookmarksPage: Int = 1,
    val followingPage: Int = 1,
    val followersPage: Int = 1,
    val hasMorePosts: Boolean = true,
    val hasMoreBookmarks: Boolean = true,
    val hasMoreFollowing: Boolean = true,
    val hasMoreFollowers: Boolean = true
)

/**
 * 个人中心 ViewModel
 */
@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val profileRepository: ProfileRepository,
    private val tokenManager: TokenManager,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    // 从导航参数获取 userId，为 null 或 0 表示查看自己的主页
    private val targetUserId: Long? = savedStateHandle.get<Long>("userId")

    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()

    init {
        checkLoginStatus()
    }

    /**
     * 检查登录状态
     */
    private fun checkLoginStatus() {
        viewModelScope.launch {
            val token = tokenManager.getToken()
            val isLoggedIn = !token.isNullOrBlank()
            _uiState.value = _uiState.value.copy(isLoggedIn = isLoggedIn)
            if (isLoggedIn) {
                if (targetUserId != null && targetUserId > 0) {
                    // 查看他人主页
                    _uiState.value = _uiState.value.copy(isOwnProfile = false)
                    loadUserProfile(targetUserId)
                } else {
                    // 查看自己的主页
                    _uiState.value = _uiState.value.copy(isOwnProfile = true)
                    loadProfile()
                }
            }
        }
    }

    /**
     * 加载当前用户信息（自己的主页）
     */
    fun loadProfile() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val response = profileRepository.getCurrentUser()

                if (response.code == 0) {
                    val user = response.data
                    _uiState.value = _uiState.value.copy(
                        user = user,
                        isLoading = false
                    )
                    if (user != null) {
                        loadUserPosts(user.id)
                    }
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "加载失败"
                )
            }
        }
    }

    /**
     * 加载指定用户信息（他人主页）
     */
    fun loadUserProfile(userId: Long) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val response = profileRepository.getUser(userId)

                if (response.code == 0) {
                    val user = response.data
                    _uiState.value = _uiState.value.copy(
                        user = user,
                        isLoading = false,
                        isFollowing = user?.isFollowing ?: false
                    )
                    if (user != null) {
                        loadUserPosts(user.id)
                    }
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "加载失败"
                )
            }
        }
    }

    /**
     * 关注/取消关注
     */
    fun toggleFollow() {
        val user = _uiState.value.user ?: return
        val isFollowing = _uiState.value.isFollowing

        viewModelScope.launch {
            try {
                val response = if (isFollowing) {
                    profileRepository.unfollowUser(user.id)
                } else {
                    profileRepository.followUser(user.id)
                }

                if (response.code == 0) {
                    _uiState.value = _uiState.value.copy(
                        isFollowing = !isFollowing,
                        user = user.copy(
                            followerCount = if (isFollowing) {
                                (user.followerCount - 1).coerceAtLeast(0)
                            } else {
                                user.followerCount + 1
                            }
                        )
                    )
                }
            } catch (_: Exception) {
                // 静默失败
            }
        }
    }

    /**
     * 加载用户发布的帖子
     */
    fun loadUserPosts(userId: Long) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoadingPosts = true)

            try {
                val response = profileRepository.getUserPosts(userId = userId, page = 1)

                if (response.code == 0) {
                    val posts = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        posts = posts,
                        isLoadingPosts = false,
                        postsPage = 1,
                        hasMorePosts = posts.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingPosts = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingPosts = false)
            }
        }
    }

    /**
     * 加载更多帖子
     */
    fun loadMorePosts() {
        val state = _uiState.value
        val userId = state.user?.id ?: return
        if (state.isLoadingPosts || !state.hasMorePosts) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoadingPosts = true)
            try {
                val nextPage = state.postsPage + 1
                val response = profileRepository.getUserPosts(userId = userId, page = nextPage)
                if (response.code == 0) {
                    val newPosts = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        posts = state.posts + newPosts,
                        isLoadingPosts = false,
                        postsPage = nextPage,
                        hasMorePosts = newPosts.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingPosts = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingPosts = false)
            }
        }
    }

    /**
     * 加载收藏列表
     */
    fun loadBookmarks() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoadingBookmarks = true)
            try {
                val response = profileRepository.getBookmarks(page = 1)
                if (response.code == 0) {
                    val bookmarks = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        bookmarks = bookmarks,
                        isLoadingBookmarks = false,
                        bookmarksPage = 1,
                        hasMoreBookmarks = bookmarks.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingBookmarks = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingBookmarks = false)
            }
        }
    }

    /**
     * 加载更多收藏
     */
    fun loadMoreBookmarks() {
        val state = _uiState.value
        if (state.isLoadingBookmarks || !state.hasMoreBookmarks) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoadingBookmarks = true)
            try {
                val nextPage = state.bookmarksPage + 1
                val response = profileRepository.getBookmarks(page = nextPage)
                if (response.code == 0) {
                    val newBookmarks = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        bookmarks = state.bookmarks + newBookmarks,
                        isLoadingBookmarks = false,
                        bookmarksPage = nextPage,
                        hasMoreBookmarks = newBookmarks.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingBookmarks = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingBookmarks = false)
            }
        }
    }

    /**
     * 加载关注列表
     */
    fun loadFollowing() {
        val userId = _uiState.value.user?.id ?: return
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoadingFollowing = true)
            try {
                val response = profileRepository.getFollowing(userId, page = 1)
                if (response.code == 0) {
                    val list = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        followingList = list,
                        isLoadingFollowing = false,
                        followingPage = 1,
                        hasMoreFollowing = list.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingFollowing = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingFollowing = false)
            }
        }
    }

    /**
     * 加载更多关注
     */
    fun loadMoreFollowing() {
        val state = _uiState.value
        val userId = state.user?.id ?: return
        if (state.isLoadingFollowing || !state.hasMoreFollowing) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoadingFollowing = true)
            try {
                val nextPage = state.followingPage + 1
                val response = profileRepository.getFollowing(userId, page = nextPage)
                if (response.code == 0) {
                    val newList = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        followingList = state.followingList + newList,
                        isLoadingFollowing = false,
                        followingPage = nextPage,
                        hasMoreFollowing = newList.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingFollowing = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingFollowing = false)
            }
        }
    }

    /**
     * 加载粉丝列表
     */
    fun loadFollowers() {
        val userId = _uiState.value.user?.id ?: return
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoadingFollowers = true)
            try {
                val response = profileRepository.getFollowers(userId, page = 1)
                if (response.code == 0) {
                    val list = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        followersList = list,
                        isLoadingFollowers = false,
                        followersPage = 1,
                        hasMoreFollowers = list.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingFollowers = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingFollowers = false)
            }
        }
    }

    /**
     * 加载更多粉丝
     */
    fun loadMoreFollowers() {
        val state = _uiState.value
        val userId = state.user?.id ?: return
        if (state.isLoadingFollowers || !state.hasMoreFollowers) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoadingFollowers = true)
            try {
                val nextPage = state.followersPage + 1
                val response = profileRepository.getFollowers(userId, page = nextPage)
                if (response.code == 0) {
                    val newList = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        followersList = state.followersList + newList,
                        isLoadingFollowers = false,
                        followersPage = nextPage,
                        hasMoreFollowers = newList.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(isLoadingFollowers = false)
                }
            } catch (_: Exception) {
                _uiState.value = _uiState.value.copy(isLoadingFollowers = false)
            }
        }
    }

    /**
     * 切换 Tab
     */
    fun selectTab(tab: ProfileTab) {
        _uiState.value = _uiState.value.copy(selectedTab = tab)
        when (tab) {
            ProfileTab.POSTS -> {
                if (_uiState.value.posts.isEmpty()) {
                    _uiState.value.user?.let { loadUserPosts(it.id) }
                }
            }
            ProfileTab.BOOKMARKS -> {
                if (_uiState.value.bookmarks.isEmpty() && _uiState.value.isOwnProfile) {
                    loadBookmarks()
                }
            }
            ProfileTab.FOLLOWING -> {
                if (_uiState.value.followingList.isEmpty()) {
                    loadFollowing()
                }
            }
            ProfileTab.FOLLOWERS -> {
                if (_uiState.value.followersList.isEmpty()) {
                    loadFollowers()
                }
            }
        }
    }

    /**
     * 退出登录
     */
    fun logout() {
        viewModelScope.launch {
            tokenManager.clearToken()
            _uiState.value = ProfileUiState(isLoggedIn = false)
        }
    }

    /**
     * 清除错误信息
     */
    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }

    companion object {
        private const val PAGE_SIZE = 20
    }
}
