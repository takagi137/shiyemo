package com.shiyemo.app.ui.post

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.model.dto.CommentDto
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.repository.PostRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

/**
 * 帖子详情 UI 状态
 */
data class PostDetailUiState(
    val post: PostDto? = null,
    val comments: List<CommentDto> = emptyList(),
    val isLoading: Boolean = false,
    val isLoadingComments: Boolean = false,
    val isSubmittingComment: Boolean = false,
    val error: String? = null,
    val commentError: String? = null,
    val currentPage: Int = 1,
    val hasMoreComments: Boolean = true,
    val commentContent: String = ""
)

/**
 * 帖子详情 ViewModel
 */
@HiltViewModel
class PostDetailViewModel @Inject constructor(
    private val postRepository: PostRepository,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    private val postId: Long = savedStateHandle.get<Long>("postId") ?: 0L

    private val _uiState = MutableStateFlow(PostDetailUiState())
    val uiState: StateFlow<PostDetailUiState> = _uiState.asStateFlow()

    init {
        if (postId > 0) {
            loadPost()
            loadComments()
        }
    }

    /**
     * 加载帖子详情
     */
    fun loadPost() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val response = postRepository.getPostDetail(postId)

                if (response.code == 0) {
                    _uiState.value = _uiState.value.copy(
                        post = response.data,
                        isLoading = false
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
     * 加载评论列表
     */
    fun loadComments() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(
                isLoadingComments = true,
                currentPage = 1
            )

            try {
                val response = postRepository.getComments(
                    postId = postId,
                    page = 1
                )

                if (response.code == 0) {
                    val comments = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        comments = comments,
                        isLoadingComments = false,
                        currentPage = 1,
                        hasMoreComments = comments.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoadingComments = false,
                        commentError = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoadingComments = false,
                    commentError = e.message ?: "加载评论失败"
                )
            }
        }
    }

    /**
     * 加载更多评论
     */
    fun loadMoreComments() {
        val state = _uiState.value
        if (state.isLoadingComments || !state.hasMoreComments) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoadingComments = true)

            try {
                val nextPage = state.currentPage + 1
                val response = postRepository.getComments(
                    postId = postId,
                    page = nextPage
                )

                if (response.code == 0) {
                    val newComments = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        comments = state.comments + newComments,
                        isLoadingComments = false,
                        currentPage = nextPage,
                        hasMoreComments = newComments.size >= PAGE_SIZE
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoadingComments = false,
                        commentError = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoadingComments = false,
                    commentError = e.message ?: "加载更多评论失败"
                )
            }
        }
    }

    /**
     * 更新评论输入内容
     */
    fun updateCommentContent(content: String) {
        _uiState.value = _uiState.value.copy(commentContent = content)
    }

    /**
     * 提交评论
     */
    fun submitComment() {
        val content = _uiState.value.commentContent.trim()
        if (content.isBlank()) return

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isSubmittingComment = true, commentError = null)

            try {
                val response = postRepository.createComment(postId, content)

                if (response.code == 0) {
                    // 添加新评论到列表头部
                    val newComment = response.data
                    if (newComment != null) {
                        _uiState.value = _uiState.value.copy(
                            comments = listOf(newComment) + _uiState.value.comments,
                            commentContent = "",
                            isSubmittingComment = false
                        )
                    }

                    // 更新帖子的评论数
                    _uiState.value.post?.let { post ->
                        _uiState.value = _uiState.value.copy(
                            post = post.copy(commentCount = post.commentCount + 1)
                        )
                    }
                } else {
                    _uiState.value = _uiState.value.copy(
                        isSubmittingComment = false,
                        commentError = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isSubmittingComment = false,
                    commentError = e.message ?: "评论失败"
                )
            }
        }
    }

    /**
     * 点赞/取消点赞
     */
    fun toggleLike() {
        viewModelScope.launch {
            try {
                val response = postRepository.toggleLike(postId)
                if (response.code == 0) {
                    _uiState.value.post?.let { post ->
                        _uiState.value = _uiState.value.copy(
                            post = post.copy(
                                isLiked = !post.isLiked,
                                likeCount = if (post.isLiked) post.likeCount - 1 else post.likeCount + 1
                            )
                        )
                    }
                }
            } catch (_: Exception) {
                // 静默失败
            }
        }
    }

    /**
     * 收藏/取消收藏
     */
    fun toggleBookmark() {
        viewModelScope.launch {
            try {
                val response = postRepository.toggleBookmark(postId)
                if (response.code == 0) {
                    _uiState.value.post?.let { post ->
                        _uiState.value = _uiState.value.copy(
                            post = post.copy(
                                isBookmarked = !post.isBookmarked,
                                bookmarkCount = if (post.isBookmarked) post.bookmarkCount - 1 else post.bookmarkCount + 1
                            )
                        )
                    }
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
        _uiState.value = _uiState.value.copy(error = null, commentError = null)
    }

    companion object {
        private const val PAGE_SIZE = 20
    }
}
