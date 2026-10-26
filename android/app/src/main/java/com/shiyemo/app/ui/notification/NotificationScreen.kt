package com.shiyemo.app.ui.notification

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.PersonAdd
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.shiyemo.app.data.model.dto.NotificationDto
import com.shiyemo.app.ui.components.EmptyState
import com.shiyemo.app.ui.components.ErrorState
import com.shiyemo.app.ui.components.LoadMoreIndicator
import com.shiyemo.app.ui.components.LoadingState
import com.shiyemo.app.util.TimeUtils

/**
 * 通知页面
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NotificationScreen(
    onNavigateToPost: (Long) -> Unit,
    onNavigateToUser: (Long) -> Unit,
    viewModel: NotificationViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val listState = rememberLazyListState()

    // 检测是否需要加载更多
    val shouldLoadMore by remember {
        derivedStateOf {
            val lastVisibleItem = listState.layoutInfo.visibleItemsInfo.lastOrNull()?.index ?: 0
            lastVisibleItem >= uiState.notifications.size - 3 && !uiState.isLoading && uiState.hasMore
        }
    }

    // 触发加载更多
    LaunchedEffect(shouldLoadMore) {
        if (shouldLoadMore) {
            viewModel.loadMore()
        }
    }

    // 未读数
    val unreadCount = uiState.notifications.count { !it.isRead }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(if (unreadCount > 0) "通知 ($unreadCount)" else "通知")
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary
                ),
                actions = {
                    if (unreadCount > 0) {
                        TextButton(onClick = { viewModel.markAllAsRead() }) {
                            Text(
                                text = "全部已读",
                                color = MaterialTheme.colorScheme.onPrimary
                            )
                        }
                    }
                }
            )
        }
    ) { paddingValues ->
        when {
            uiState.isLoading && uiState.notifications.isEmpty() -> {
                LoadingState(modifier = Modifier.padding(paddingValues))
            }
            uiState.error != null && uiState.notifications.isEmpty() -> {
                ErrorState(
                    message = uiState.error ?: "加载失败",
                    onRetry = { viewModel.loadNotifications() },
                    modifier = Modifier.padding(paddingValues)
                )
            }
            uiState.notifications.isEmpty() -> {
                EmptyState(
                    icon = Icons.Default.Notifications,
                    message = "暂无通知",
                    modifier = Modifier.padding(paddingValues)
                )
            }
            else -> {
                LazyColumn(
                    state = listState,
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(paddingValues),
                    verticalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    items(
                        items = uiState.notifications,
                        key = { it.id }
                    ) { notification ->
                        NotificationItem(
                            notification = notification,
                            onClick = {
                                // 标记为已读
                                if (!notification.isRead) {
                                    viewModel.markAsRead(notification.id)
                                }
                                // 导航
                                when (notification.type) {
                                    "like", "comment" -> {
                                        if (notification.postId > 0) {
                                            onNavigateToPost(notification.postId)
                                        }
                                    }
                                    "follow" -> {
                                        notification.fromUser?.let {
                                            onNavigateToUser(it.id)
                                        }
                                    }
                                }
                            }
                        )
                    }

                    // 加载更多指示器
                    if (uiState.isLoading) {
                        item { LoadMoreIndicator() }
                    }

                    // 底部间距
                    item {
                        Spacer(modifier = Modifier.height(16.dp))
                    }
                }
            }
        }
    }
}

/**
 * 通知项
 */
@Composable
private fun NotificationItem(
    notification: NotificationDto,
    onClick: () -> Unit
) {
    val backgroundColor = if (notification.isRead) {
        MaterialTheme.colorScheme.surface
    } else {
        MaterialTheme.colorScheme.primaryContainer.copy(alpha = 0.3f)
    }

    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 2.dp)
            .clickable(onClick = onClick),
        colors = CardDefaults.cardColors(containerColor = backgroundColor),
        elevation = CardDefaults.cardElevation(defaultElevation = 0.dp)
    ) {
        Row(
            modifier = Modifier.padding(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // 通知图标
            val (icon, iconTint) = when (notification.type) {
                "like" -> Icons.Default.Favorite to Color.Red
                "comment" -> Icons.Default.CheckCircle to MaterialTheme.colorScheme.primary
                "follow" -> Icons.Default.PersonAdd to MaterialTheme.colorScheme.tertiary
                else -> Icons.Default.Notifications to MaterialTheme.colorScheme.onSurfaceVariant
            }

            Surface(
                shape = CircleShape,
                color = iconTint.copy(alpha = 0.1f),
                modifier = Modifier.size(40.dp)
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Icon(
                        imageVector = icon,
                        contentDescription = null,
                        modifier = Modifier.size(24.dp),
                        tint = iconTint
                    )
                }
            }

            Spacer(modifier = Modifier.width(12.dp))

            Column(modifier = Modifier.weight(1f)) {
                // 通知内容
                Text(
                    text = notification.content.ifBlank {
                        when (notification.type) {
                            "like" -> "赞了你的帖子"
                            "comment" -> "评论了你的帖子"
                            "follow" -> "关注了你"
                            else -> "新通知"
                        }
                    },
                    style = MaterialTheme.typography.bodyMedium,
                    maxLines = 2,
                    overflow = TextOverflow.Ellipsis,
                    fontWeight = if (notification.isRead) FontWeight.Normal else FontWeight.Bold
                )

                // 来自谁 + 时间
                Row(
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    notification.fromUser?.let { user ->
                        Text(
                            text = user.nickname.ifBlank { "用户${user.id}" },
                            style = MaterialTheme.typography.bodySmall,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                        Text(
                            text = " · ",
                            style = MaterialTheme.typography.bodySmall,
                            color = MaterialTheme.colorScheme.onSurfaceVariant.copy(alpha = 0.5f)
                        )
                    }
                    Text(
                        text = TimeUtils.formatRelative(notification.createdAt),
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            // 未读标记
            if (!notification.isRead) {
                Surface(
                    shape = CircleShape,
                    color = MaterialTheme.colorScheme.primary,
                    modifier = Modifier.size(8.dp)
                ) {}
            }
        }
    }
}
