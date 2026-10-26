package com.shiyemo.app.ui.post

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.FlowRow
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
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.Bookmark
import androidx.compose.material.icons.filled.Comment
import androidx.compose.material.icons.filled.Favorite
import androidx.compose.material.icons.filled.Send
import androidx.compose.material.icons.outlined.BookmarkBorder
import androidx.compose.material.icons.outlined.FavoriteBorder
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
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
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import coil.compose.AsyncImage
import com.shiyemo.app.data.model.dto.CommentDto
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.ui.components.EmptyState
import com.shiyemo.app.ui.components.ErrorState
import com.shiyemo.app.ui.components.LoadMoreIndicator
import com.shiyemo.app.ui.components.LoadingState
import com.shiyemo.app.util.TimeUtils

/**
 * 帖子详情页面
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PostDetailScreen(
    onBack: () -> Unit,
    viewModel: PostDetailViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val listState = rememberLazyListState()

    // 检测是否需要加载更多评论
    val shouldLoadMoreComments by remember {
        derivedStateOf {
            val lastVisibleItem = listState.layoutInfo.visibleItemsInfo.lastOrNull()?.index ?: 0
            lastVisibleItem >= uiState.comments.size + 1 && !uiState.isLoadingComments && uiState.hasMoreComments
        }
    }

    // 触发加载更多评论
    LaunchedEffect(shouldLoadMoreComments) {
        if (shouldLoadMoreComments) {
            viewModel.loadMoreComments()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("帖子详情") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(
                            imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                            contentDescription = "返回"
                        )
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary,
                    navigationIconContentColor = MaterialTheme.colorScheme.onPrimary
                )
            )
        }
    ) { paddingValues ->
        if (uiState.isLoading && uiState.post == null) {
            LoadingState(modifier = Modifier.padding(paddingValues))
        } else if (uiState.error != null && uiState.post == null) {
            ErrorState(
                message = uiState.error ?: "加载失败",
                onRetry = { viewModel.loadPost() },
                modifier = Modifier.padding(paddingValues)
            )
        } else if (uiState.post != null) {
            // 帖子详情
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(paddingValues)
            ) {
                // 帖子内容区域（可滚动）
                LazyColumn(
                    state = listState,
                    modifier = Modifier.weight(1f)
                ) {
                    // 帖子详情
                    item {
                        PostDetailContent(
                            post = uiState.post!!,
                            onLikeClick = { viewModel.toggleLike() },
                            onBookmarkClick = { viewModel.toggleBookmark() }
                        )
                    }

                    // 评论区标题
                    item {
                        Text(
                            text = "评论 (${uiState.post!!.commentCount})",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold,
                            modifier = Modifier.padding(16.dp)
                        )
                    }

                    // 评论列表
                    if (uiState.isLoadingComments && uiState.comments.isEmpty()) {
                        item { LoadMoreIndicator() }
                    } else if (uiState.comments.isEmpty()) {
                        item {
                            EmptyState(
                                icon = Icons.Default.Comment,
                                message = "暂无评论",
                                subtitle = "快来抢沙发吧~",
                                modifier = Modifier.height(200.dp)
                            )
                        }
                    } else {
                        items(
                            items = uiState.comments,
                            key = { it.id }
                        ) { comment ->
                            CommentItem(comment = comment)
                        }

                        // 加载更多评论指示器
                        if (uiState.isLoadingComments) {
                            item { LoadMoreIndicator() }
                        }
                    }

                    // 底部间距
                    item {
                        Spacer(modifier = Modifier.height(16.dp))
                    }
                }

                // 评论输入框
                CommentInput(
                    content = uiState.commentContent,
                    onContentChange = { viewModel.updateCommentContent(it) },
                    onSubmit = { viewModel.submitComment() },
                    isSubmitting = uiState.isSubmittingComment,
                    error = uiState.commentError
                )
            }
        }
    }
}

/**
 * 帖子详情内容
 */
@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun PostDetailContent(
    post: PostDto,
    onLikeClick: () -> Unit,
    onBookmarkClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            // 用户信息
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    imageVector = Icons.Default.AccountCircle,
                    contentDescription = null,
                    modifier = Modifier.size(48.dp),
                    tint = MaterialTheme.colorScheme.primary
                )
                Spacer(modifier = Modifier.width(12.dp))
                Column {
                    Text(
                        text = if (post.isAnonymous) post.anonymousName else (post.user?.nickname ?: "用户${post.userId}"),
                        style = MaterialTheme.typography.titleMedium,
                        fontWeight = FontWeight.Bold
                    )
                    Text(
                        text = formatTime(post.createdAt),
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            Spacer(modifier = Modifier.height(16.dp))

            // 帖子标题
            if (post.title.isNotBlank()) {
                Text(
                    text = post.title,
                    style = MaterialTheme.typography.headlineSmall,
                    fontWeight = FontWeight.Bold
                )
                Spacer(modifier = Modifier.height(12.dp))
            }

            // 帖子内容
            Text(
                text = post.content,
                style = MaterialTheme.typography.bodyLarge
            )

            // 图片列表
            if (post.images.isNotEmpty()) {
                Spacer(modifier = Modifier.height(16.dp))
                post.images.forEach { imageUrl ->
                    AsyncImage(
                        model = imageUrl,
                        contentDescription = null,
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(300.dp)
                            .clip(MaterialTheme.shapes.medium),
                        contentScale = ContentScale.Crop
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }
            }

            // 标签
            if (post.tags.isNotEmpty()) {
                Spacer(modifier = Modifier.height(16.dp))
                FlowRow(
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    post.tags.forEach { tag ->
                        Text(
                            text = "#$tag",
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.primary
                        )
                    }
                }
            }

            Spacer(modifier = Modifier.height(16.dp))

            // 互动按钮
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                // 点赞
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.padding(8.dp)
                ) {
                    IconButton(onClick = onLikeClick) {
                        Icon(
                            imageVector = if (post.isLiked) Icons.Default.Favorite else Icons.Outlined.FavoriteBorder,
                            contentDescription = "点赞",
                            tint = if (post.isLiked) Color.Red else MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                    Text(
                        text = "${post.likeCount}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }

                // 评论
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.padding(8.dp)
                ) {
                    Icon(
                        imageVector = Icons.Default.Comment,
                        contentDescription = "评论",
                        modifier = Modifier.size(24.dp),
                        tint = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "${post.commentCount}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }

                // 收藏
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    modifier = Modifier.padding(8.dp)
                ) {
                    IconButton(onClick = onBookmarkClick) {
                        Icon(
                            imageVector = if (post.isBookmarked) Icons.Default.Bookmark else Icons.Outlined.BookmarkBorder,
                            contentDescription = "收藏",
                            tint = if (post.isBookmarked) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                    Text(
                        text = "${post.bookmarkCount}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
            }
        }
    }
}

/**
 * 评论项
 */
@Composable
private fun CommentItem(comment: CommentDto) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 4.dp),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Column(
            modifier = Modifier.padding(12.dp)
        ) {
            // 用户信息
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    imageVector = Icons.Default.AccountCircle,
                    contentDescription = null,
                    modifier = Modifier.size(32.dp),
                    tint = MaterialTheme.colorScheme.primary
                )
                Spacer(modifier = Modifier.width(8.dp))
                Column {
                    Text(
                        text = if (comment.isAnonymous) comment.anonymousName else (comment.user?.nickname ?: "用户${comment.userId}"),
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.Bold
                    )
                    Text(
                        text = formatTime(comment.createdAt),
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            Spacer(modifier = Modifier.height(8.dp))

            // 评论内容
            Text(
                text = comment.content,
                style = MaterialTheme.typography.bodyMedium
            )
        }
    }
}

/**
 * 评论输入框
 */
@Composable
private fun CommentInput(
    content: String,
    onContentChange: (String) -> Unit,
    onSubmit: () -> Unit,
    isSubmitting: Boolean,
    error: String?
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        elevation = CardDefaults.cardElevation(defaultElevation = 8.dp)
    ) {
        Column(
            modifier = Modifier.padding(12.dp)
        ) {
            // 错误提示
            if (error != null) {
                Text(
                    text = error,
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.error,
                    modifier = Modifier.padding(bottom = 8.dp)
                )
            }

            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                OutlinedTextField(
                    value = content,
                    onValueChange = onContentChange,
                    placeholder = { Text("写评论...") },
                    modifier = Modifier.weight(1f),
                    maxLines = 3,
                    enabled = !isSubmitting
                )

                Spacer(modifier = Modifier.width(8.dp))

                IconButton(
                    onClick = onSubmit,
                    enabled = content.isNotBlank() && !isSubmitting
                ) {
                    if (isSubmitting) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(24.dp)
                        )
                    } else {
                        Icon(
                            imageVector = Icons.Default.Send,
                            contentDescription = "发送",
                            tint = if (content.isNotBlank()) {
                                MaterialTheme.colorScheme.primary
                            } else {
                                MaterialTheme.colorScheme.onSurfaceVariant
                            }
                        )
                    }
                }
            }
        }
    }
}

private fun formatTime(createdAt: String): String = TimeUtils.formatRelative(createdAt)
