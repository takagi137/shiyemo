package com.shiyemo.app.ui.feed

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.repository.FeedRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

/**
 * Feed 页面 UI 状态
 */
data class FeedUiState(
    val posts: List<PostDto> = emptyList(),
    val isLoading: Boolean = false,
    val isLoadingMore: Boolean = false,
    val error: String? = null,
    val currentPage: Int = 1,
    val hasMore: Boolean = true,
    val selectedType: String? = null // null 表示全部
)

/**
 * Feed ViewModel
 */
@HiltViewModel
class FeedViewModel @Inject constructor(
    private val feedRepository: FeedRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(FeedUiState())
    val uiState: StateFlow<FeedUiState> = _uiState.asStateFlow()

    init {
        loadPosts()
    }

    /**
     * 加载帖子列表（首次加载或刷新）
     */
    fun loadPosts() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(
                isLoading = true,
                error = null,
                currentPage = 1
            )

            try {
                val response = feedRepository.getPosts(
                    type = _uiState.value.selectedType,
                    page = 1,
                    pageSize = PAGE_SIZE
                )

                if (response.code == 0) {
                    val posts = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        posts = posts,
                        isLoading = false,
                        currentPage = 1,
                        hasMore = posts.size >= PAGE_SIZE
                    )
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
     * 加载更多（分页）
     */
    fun loadMore() {
        val state = _uiState.value
        if (state.isLoadingMore || !state.hasMore || state.isLoading) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoadingMore = true)

            try {
                val nextPage = state.currentPage + 1
                val response = feedRepository.getPosts(
                    type = state.selectedType,
                    page = nextPage,
                    pageSize = PAGE_SIZE
                )

                if (response.code == 0) {
                    val newPosts = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        posts = state.posts + newPosts,
                        isLoadingMore = false,
                        currentPage = nextPage,
                        hasMore = newPosts.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoadingMore = false,
                        error = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoadingMore = false,
                    error = e.message ?: "加载更多失败"
                )
            }
        }
    }

    /**
     * 切换帖子类型过滤
     */
    fun selectType(type: String?) {
        if (_uiState.value.selectedType == type) return
        _uiState.value = _uiState.value.copy(selectedType = type)
        loadPosts()
    }

    /**
     * 点赞/取消点赞
     */
    fun toggleLike(postId: Long) {
        viewModelScope.launch {
            try {
                val response = feedRepository.toggleLike(postId)
                if (response.code == 0) {
                    // 更新本地状态
                    val updatedPosts = _uiState.value.posts.map { post ->
                        if (post.id == postId) {
                            post.copy(
                                isLiked = !post.isLiked,
                                likeCount = if (post.isLiked) post.likeCount - 1 else post.likeCount + 1
                            )
                        } else {
                            post
                        }
                    }
                    _uiState.value = _uiState.value.copy(posts = updatedPosts)
                }
            } catch (_: Exception) {
                // 静默失败，不影响用户体验
            }
        }
    }

    /**
     * 收藏/取消收藏
     */
    fun toggleBookmark(postId: Long) {
        viewModelScope.launch {
            try {
                val response = feedRepository.toggleBookmark(postId)
                if (response.code == 0) {
                    // 更新本地状态
                    val updatedPosts = _uiState.value.posts.map { post ->
                        if (post.id == postId) {
                            post.copy(
                                isBookmarked = !post.isBookmarked,
                                bookmarkCount = if (post.isBookmarked) post.bookmarkCount - 1 else post.bookmarkCount + 1
                            )
                        } else {
                            post
                        }
                    }
                    _uiState.value = _uiState.value.copy(posts = updatedPosts)
                }
            } catch (_: Exception) {
                // 静默失败
            }
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
