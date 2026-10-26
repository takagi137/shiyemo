package com.shiyemo.app.ui.notification

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.model.dto.NotificationDto
import com.shiyemo.app.data.repository.NotificationRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

/**
 * 通知页面 UI 状态
 */
data class NotificationUiState(
    val notifications: List<NotificationDto> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val currentPage: Int = 1,
    val hasMore: Boolean = true
)

/**
 * 通知 ViewModel
 */
@HiltViewModel
class NotificationViewModel @Inject constructor(
    private val notificationRepository: NotificationRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(NotificationUiState())
    val uiState: StateFlow<NotificationUiState> = _uiState.asStateFlow()

    init {
        loadNotifications()
    }

    /**
     * 加载通知列表（首次加载或刷新）
     */
    fun loadNotifications() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(
                isLoading = true,
                error = null,
                currentPage = 1
            )

            try {
                val response = notificationRepository.getNotifications(
                    page = 1,
                    pageSize = PAGE_SIZE
                )

                if (response.code == 0) {
                    val notifications = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        notifications = notifications,
                        isLoading = false,
                        currentPage = 1,
                        hasMore = notifications.size >= PAGE_SIZE
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
                    error = e.message ?: "加载通知失败"
                )
            }
        }
    }

    /**
     * 加载更多
     */
    fun loadMore() {
        val state = _uiState.value
        if (state.isLoading || !state.hasMore) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoading = true)

            try {
                val nextPage = state.currentPage + 1
                val response = notificationRepository.getNotifications(
                    page = nextPage,
                    pageSize = PAGE_SIZE
                )

                if (response.code == 0) {
                    val newNotifications = response.data ?: emptyList()
                    _uiState.value = _uiState.value.copy(
                        notifications = state.notifications + newNotifications,
                        isLoading = false,
                        currentPage = nextPage,
                        hasMore = newNotifications.size >= PAGE_SIZE
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
                    error = e.message ?: "加载更多失败"
                )
            }
        }
    }

    /**
     * 标记单条通知为已读
     */
    fun markAsRead(notificationId: Long) {
        viewModelScope.launch {
            try {
                val response = notificationRepository.markAsRead(notificationId)
                if (response.code == 0) {
                    // 更新本地状态
                    val updated = _uiState.value.notifications.map {
                        if (it.id == notificationId) it.copy(isRead = true) else it
                    }
                    _uiState.value = _uiState.value.copy(notifications = updated)
                }
            } catch (_: Exception) {
                // 静默失败
            }
        }
    }

    /**
     * 标记所有通知为已读
     */
    fun markAllAsRead() {
        viewModelScope.launch {
            try {
                val response = notificationRepository.markAllAsRead()
                if (response.code == 0) {
                    val updated = _uiState.value.notifications.map { it.copy(isRead = true) }
                    _uiState.value = _uiState.value.copy(notifications = updated)
                }
            } catch (_: Exception) {
                // 静默失败
            }
        }
    }

    companion object {
        private const val PAGE_SIZE = 20
    }
}
