package com.shiyemo.app.ui.search

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
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
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.Clear
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.material3.TextFieldDefaults
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.shiyemo.app.data.model.dto.SearchResultDto
import com.shiyemo.app.ui.components.EmptyState
import com.shiyemo.app.ui.components.ErrorState
import com.shiyemo.app.ui.components.LoadMoreIndicator
import com.shiyemo.app.ui.components.LoadingState

/**
 * 搜索页面
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SearchScreen(
    onNavigateToPost: (Long) -> Unit,
    onNavigateToUser: (Long) -> Unit,
    viewModel: SearchViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val listState = rememberLazyListState()

    // 检测是否需要加载更多
    val shouldLoadMore by remember {
        derivedStateOf {
            val lastVisibleItem = listState.layoutInfo.visibleItemsInfo.lastOrNull()?.index ?: 0
            lastVisibleItem >= uiState.results.size - 3 && !uiState.isLoading && uiState.hasMore
        }
    }

    // 触发加载更多
    LaunchedEffect(shouldLoadMore) {
        if (shouldLoadMore && uiState.isSearchPerformed) {
            viewModel.loadMore()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("搜索") },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary
                )
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // 搜索输入框
            SearchBar(
                query = uiState.query,
                onQueryChange = { viewModel.updateQuery(it) },
                onSearch = { viewModel.search() },
                onClear = { viewModel.clearSearch() }
            )

            // 内容区域
            when {
                uiState.isLoading && uiState.results.isEmpty() -> {
                    LoadingState()
                }
                uiState.error != null && uiState.results.isEmpty() -> {
                    ErrorState(
                        message = uiState.error ?: "搜索失败",
                        onRetry = { viewModel.search() }
                    )
                }
                uiState.isSearchPerformed && uiState.results.isEmpty() && !uiState.isLoading -> {
                    EmptyState(
                        icon = Icons.Default.Search,
                        message = "未找到相关内容"
                    )
                }
                !uiState.isSearchPerformed -> {
                    EmptyState(
                        icon = Icons.Default.Search,
                        message = "搜索帖子或用户",
                        subtitle = "输入关键词开始搜索"
                    )
                }
                else -> {
                    // 搜索结果列表
                    LazyColumn(
                        state = listState,
                        modifier = Modifier.fillMaxSize(),
                        verticalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        item {
                            Text(
                                text = "搜索结果 (${uiState.results.size})",
                                style = MaterialTheme.typography.titleSmall,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                                modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
                            )
                        }

                        items(
                            items = uiState.results,
                            key = { "${it.type}-${it.id}" }
                        ) { result ->
                            SearchResultCard(
                                result = result,
                                onClick = {
                                    when (result.type) {
                                        "post" -> onNavigateToPost(result.id)
                                        "user" -> onNavigateToUser(result.id)
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
}

/**
 * 搜索栏
 */
@Composable
private fun SearchBar(
    query: String,
    onQueryChange: (String) -> Unit,
    onSearch: () -> Unit,
    onClear: () -> Unit
) {
    TextField(
        value = query,
        onValueChange = onQueryChange,
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 8.dp),
        placeholder = { Text("搜索帖子或用户...") },
        leadingIcon = {
            Icon(
                imageVector = Icons.Default.Search,
                contentDescription = "搜索"
            )
        },
        trailingIcon = {
            if (query.isNotBlank()) {
                IconButton(onClick = onClear) {
                    Icon(
                        imageVector = Icons.Default.Clear,
                        contentDescription = "清除"
                    )
                }
            }
        },
        singleLine = true,
        colors = TextFieldDefaults.colors(
            focusedIndicatorColor = Color.Transparent,
            unfocusedIndicatorColor = Color.Transparent
        )
    )
}

/**
 * 搜索结果卡片
 */
@Composable
private fun SearchResultCard(
    result: SearchResultDto,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp)
            .clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // 图标
            Icon(
                imageVector = Icons.Default.AccountCircle,
                contentDescription = null,
                modifier = Modifier.size(40.dp),
                tint = when (result.type) {
                    "user" -> MaterialTheme.colorScheme.tertiary
                    else -> MaterialTheme.colorScheme.primary
                }
            )

            Spacer(modifier = Modifier.width(12.dp))

            Column(modifier = Modifier.weight(1f)) {
                // 类型标签
                Text(
                    text = when (result.type) {
                        "post" -> "帖子"
                        "user" -> "用户"
                        else -> result.type
                    },
                    style = MaterialTheme.typography.labelSmall,
                    color = MaterialTheme.colorScheme.primary
                )

                // 标题或昵称
                Text(
                    text = when (result.type) {
                        "user" -> result.nickname.ifBlank { "用户${result.id}" }
                        else -> result.title.ifBlank { result.content.take(50) }
                    },
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.Bold,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )

                // 内容预览
                if (result.type == "post" && result.content.isNotBlank()) {
                    Text(
                        text = result.content,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        maxLines = 2,
                        overflow = TextOverflow.Ellipsis
                    )
                }

                if (result.type == "user" && result.content.isNotBlank()) {
                    Text(
                        text = result.content,
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
