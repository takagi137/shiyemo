package com.shiyemo.app.ui.profile

import androidx.compose.foundation.clickable
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
import androidx.compose.material.icons.filled.LocationOn
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Work
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Tab
import androidx.compose.material3.TabRow
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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.shiyemo.app.data.model.dto.PostDto
import com.shiyemo.app.data.model.dto.UserDto
import com.shiyemo.app.ui.components.EmptyState
import com.shiyemo.app.ui.components.ErrorState
import com.shiyemo.app.ui.components.LoadingState

/**
 * 个人中心页面
 * @param showBackButton 是否显示返回按钮（他人主页时为 true）
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProfileScreen(
    onNavigateToLogin: () -> Unit,
    onNavigateToSettings: () -> Unit,
    onNavigateToPost: (Long) -> Unit = {},
    onNavigateToUser: (Long) -> Unit = {},
    onBack: () -> Unit = {},
    showBackButton: Boolean = false,
    viewModel: ProfileViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()

    // 未登录时跳转登录页（仅自己的主页）
    LaunchedEffect(uiState.isLoggedIn, uiState.isOwnProfile) {
        if (!uiState.isLoggedIn && uiState.isOwnProfile) {
            onNavigateToLogin()
        }
    }

    // Tab 定义（根据是否是自己的主页决定显示哪些）
    val tabs = if (uiState.isOwnProfile) {
        listOf(ProfileTab.POSTS, ProfileTab.BOOKMARKS, ProfileTab.FOLLOWING, ProfileTab.FOLLOWERS)
    } else {
        listOf(ProfileTab.POSTS, ProfileTab.FOLLOWING, ProfileTab.FOLLOWERS)
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(if (uiState.isOwnProfile) "我的" else uiState.user?.nickname ?: "用户主页")
                },
                navigationIcon = {
                    if (showBackButton) {
                        IconButton(onClick = onBack) {
                            Icon(
                                imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                                contentDescription = "返回"
                            )
                        }
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary,
                    navigationIconContentColor = MaterialTheme.colorScheme.onPrimary
                ),
                actions = {
                    if (uiState.isOwnProfile) {
                        OutlinedButton(
                            onClick = { viewModel.logout() },
                            colors = ButtonDefaults.outlinedButtonColors(
                                contentColor = MaterialTheme.colorScheme.onPrimary
                            )
                        ) {
                            Text("退出登录")
                        }
                    }
                }
            )
        }
    ) { paddingValues ->
        if (uiState.isLoading && uiState.user == null) {
            LoadingState(modifier = Modifier.padding(paddingValues))
        } else if (uiState.error != null && uiState.user == null) {
            ErrorState(
                message = uiState.error ?: "加载失败",
                onRetry = { viewModel.loadProfile() },
                modifier = Modifier.padding(paddingValues)
            )
        } else if (uiState.user != null) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(paddingValues)
            ) {
                // 用户信息卡片（带关注按钮）
                UserInfoCard(
                    user = uiState.user!!,
                    isOwnProfile = uiState.isOwnProfile,
                    isFollowing = uiState.isFollowing,
                    onFollowClick = { viewModel.toggleFollow() }
                )

                // Tab 切换
                val selectedTabIndex = tabs.indexOf(uiState.selectedTab).coerceAtLeast(0)
                TabRow(selectedTabIndex = selectedTabIndex) {
                    tabs.forEachIndexed { index, tab ->
                        Tab(
                            selected = uiState.selectedTab == tab,
                            onClick = { viewModel.selectTab(tab) },
                            text = { Text(tab.title) }
                        )
                    }
                }

                // 内容区域
                when (uiState.selectedTab) {
                    ProfileTab.POSTS -> {
                        PostsContent(
                            posts = uiState.posts,
                            isLoading = uiState.isLoadingPosts,
                            hasMore = uiState.hasMorePosts,
                            onLoadMore = { viewModel.loadMorePosts() },
                            onPostClick = onNavigateToPost
                        )
                    }
                    ProfileTab.BOOKMARKS -> {
                        if (uiState.isOwnProfile) {
                            PostsContent(
                                posts = uiState.bookmarks,
                                isLoading = uiState.isLoadingBookmarks,
                                hasMore = uiState.hasMoreBookmarks,
                                onLoadMore = { viewModel.loadMoreBookmarks() },
                                onPostClick = onNavigateToPost
                            )
                        }
                    }
                    ProfileTab.FOLLOWING -> {
                        UserListContent(
                            users = uiState.followingList,
                            isLoading = uiState.isLoadingFollowing,
                            hasMore = uiState.hasMoreFollowing,
                            onLoadMore = { viewModel.loadMoreFollowing() },
                            onUserClick = onNavigateToUser
                        )
                    }
                    ProfileTab.FOLLOWERS -> {
                        UserListContent(
                            users = uiState.followersList,
                            isLoading = uiState.isLoadingFollowers,
                            hasMore = uiState.hasMoreFollowers,
                            onLoadMore = { viewModel.loadMoreFollowers() },
                            onUserClick = onNavigateToUser
                        )
                    }
                }
            }
        }
    }
}

/**
 * 用户信息卡片（支持关注按钮）
 */
@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun UserInfoCard(
    user: UserDto,
    isOwnProfile: Boolean,
    isFollowing: Boolean,
    onFollowClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            Row(
                verticalAlignment = Alignment.CenterVertically
            ) {
                // 头像
                Icon(
                    imageVector = Icons.Default.AccountCircle,
                    contentDescription = null,
                    modifier = Modifier.size(72.dp),
                    tint = MaterialTheme.colorScheme.primary
                )

                Spacer(modifier = Modifier.width(16.dp))

                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = user.nickname.ifBlank { "用户${user.id}" },
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.Bold
                    )
                    if (user.email.isNotBlank()) {
                        Text(
                            text = user.email,
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                    }
                }

                // 关注按钮（仅他人主页显示）
                if (!isOwnProfile) {
                    Button(
                        onClick = onFollowClick,
                        colors = if (isFollowing) {
                            ButtonDefaults.outlinedButtonColors()
                        } else {
                            ButtonDefaults.buttonColors()
                        }
                    ) {
                        Text(if (isFollowing) "已关注" else "关注")
                    }
                }
            }

            Spacer(modifier = Modifier.height(16.dp))

            if (user.bio.isNotBlank()) {
                Text(text = user.bio, style = MaterialTheme.typography.bodyMedium)
                Spacer(modifier = Modifier.height(12.dp))
            }

            FlowRow(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                if (user.city.isNotBlank()) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(Icons.Default.LocationOn, null, Modifier.size(16.dp), MaterialTheme.colorScheme.onSurfaceVariant)
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(user.city, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    }
                }
                if (user.industry.isNotBlank()) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(Icons.Default.Work, null, Modifier.size(16.dp), MaterialTheme.colorScheme.onSurfaceVariant)
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(user.industry, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    }
                }
            }

            Spacer(modifier = Modifier.height(12.dp))

            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                StatItem(label = "帖子", count = user.postCount)
                StatItem(label = "关注", count = user.followingCount)
                StatItem(label = "粉丝", count = user.followerCount)
            }
        }
    }
}

@Composable
private fun StatItem(label: String, count: Int) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text("$count", style = MaterialTheme.typography.titleLarge, fontWeight = FontWeight.Bold)
        Text(label, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
    }
}

/**
 * 帖子列表内容
 */
@Composable
private fun PostsContent(
    posts: List<PostDto>,
    isLoading: Boolean,
    hasMore: Boolean,
    onLoadMore: () -> Unit,
    onPostClick: (Long) -> Unit
) {
    val listState = rememberLazyListState()

    val shouldLoadMore by remember {
        derivedStateOf {
            val lastVisibleItem = listState.layoutInfo.visibleItemsInfo.lastOrNull()?.index ?: 0
            lastVisibleItem >= posts.size - 3 && !isLoading && hasMore
        }
    }

    LaunchedEffect(shouldLoadMore) {
        if (shouldLoadMore) onLoadMore()
    }

    if (isLoading && posts.isEmpty()) {
        LoadingState()
    } else if (posts.isEmpty()) {
        EmptyState(
            icon = Icons.Default.Inbox,
            message = "暂无内容"
        )
    } else {
        LazyColumn(state = listState, modifier = Modifier.fillMaxSize(), verticalArrangement = Arrangement.spacedBy(8.dp)) {
            items(posts, key = { it.id }) { post ->
                PostItem(post = post, onClick = { onPostClick(post.id) })
            }
            if (isLoading) {
                item { LoadMoreIndicator() }
            }
            item { Spacer(modifier = Modifier.height(16.dp)) }
        }
    }
}

/**
 * 帖子项
 */
@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun PostItem(post: PostDto, onClick: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth().padding(horizontal = 16.dp).clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            if (post.title.isNotBlank()) {
                Text(post.title, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold, maxLines = 2, overflow = TextOverflow.Ellipsis)
                Spacer(modifier = Modifier.height(8.dp))
            }
            Text(post.content, style = MaterialTheme.typography.bodyMedium, maxLines = 3, overflow = TextOverflow.Ellipsis)
            Spacer(modifier = Modifier.height(12.dp))
            Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(Icons.Default.Favorite, "点赞", Modifier.size(16.dp), MaterialTheme.colorScheme.onSurfaceVariant)
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("${post.likeCount}", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(Icons.Default.Comment, "评论", Modifier.size(16.dp), MaterialTheme.colorScheme.onSurfaceVariant)
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("${post.commentCount}", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Icon(Icons.Default.Bookmark, "收藏", Modifier.size(16.dp), MaterialTheme.colorScheme.onSurfaceVariant)
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("${post.bookmarkCount}", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
            }
        }
    }
}

/**
 * 用户列表内容（关注/粉丝）
 */
@Composable
private fun UserListContent(
    users: List<UserDto>,
    isLoading: Boolean,
    hasMore: Boolean,
    onLoadMore: () -> Unit,
    onUserClick: (Long) -> Unit
) {
    val listState = rememberLazyListState()

    val shouldLoadMore by remember {
        derivedStateOf {
            val lastVisibleItem = listState.layoutInfo.visibleItemsInfo.lastOrNull()?.index ?: 0
            lastVisibleItem >= users.size - 3 && !isLoading && hasMore
        }
    }

    LaunchedEffect(shouldLoadMore) {
        if (shouldLoadMore) onLoadMore()
    }

    if (isLoading && users.isEmpty()) {
        LoadingState()
    } else if (users.isEmpty()) {
        EmptyState(
            icon = Icons.Default.Person,
            message = "暂无数据"
        )
    } else {
        LazyColumn(state = listState, modifier = Modifier.fillMaxSize(), verticalArrangement = Arrangement.spacedBy(4.dp)) {
            items(users, key = { it.id }) { user ->
                UserListItem(user = user, onClick = { onUserClick(user.id) })
            }
            if (isLoading) {
                item { LoadMoreIndicator() }
            }
            item { Spacer(modifier = Modifier.height(16.dp)) }
        }
    }
}

/**
 * 用户列表项
 */
@Composable
private fun UserListItem(user: UserDto, onClick: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth().padding(horizontal = 16.dp).clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp)
    ) {
        Row(
            modifier = Modifier.padding(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                Icons.Default.AccountCircle, null,
                modifier = Modifier.size(48.dp),
                tint = MaterialTheme.colorScheme.primary
            )
            Spacer(modifier = Modifier.width(12.dp))
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    user.nickname.ifBlank { "用户${user.id}" },
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold
                )
                if (user.bio.isNotBlank()) {
                    Text(
                        user.bio,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        maxLines = 1,
                        overflow = TextOverflow.Ellipsis
                    )
                }
            }
        }
    }
}
