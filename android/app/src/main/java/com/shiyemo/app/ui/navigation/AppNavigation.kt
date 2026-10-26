package com.shiyemo.app.ui.navigation

import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Search
import androidx.compose.material.icons.outlined.Home
import androidx.compose.material.icons.outlined.Notifications
import androidx.compose.material.icons.outlined.Person
import androidx.compose.material.icons.outlined.Search
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.shiyemo.app.ui.auth.LoginScreen
import com.shiyemo.app.ui.auth.RegisterScreen
import com.shiyemo.app.ui.feed.FeedScreen
import com.shiyemo.app.ui.notification.NotificationScreen
import com.shiyemo.app.ui.post.CreatePostScreen
import com.shiyemo.app.ui.post.PostDetailScreen
import com.shiyemo.app.ui.profile.ProfileScreen
import com.shiyemo.app.ui.search.SearchScreen

/**
 * 导航路由定义
 */
sealed class Screen(val route: String) {
    object Login : Screen("login")
    object Register : Screen("register")
    object Feed : Screen("feed")
    object Search : Screen("search")
    object CreatePost : Screen("create_post")
    object Notifications : Screen("notifications")
    object Profile : Screen("profile")
    object PostDetail : Screen("post/{postId}") {
        fun createRoute(postId: Long) = "post/$postId"
    }
    object UserProfile : Screen("user/{userId}") {
        fun createRoute(userId: Long) = "user/$userId"
    }
}

/**
 * 底部导航项
 */
data class BottomNavItem(
    val screen: Screen,
    val label: String,
    val selectedIcon: ImageVector,
    val unselectedIcon: ImageVector
)

val bottomNavItems = listOf(
    BottomNavItem(Screen.Feed, "首页", Icons.Filled.Home, Icons.Outlined.Home),
    BottomNavItem(Screen.Search, "搜索", Icons.Filled.Search, Icons.Outlined.Search),
    BottomNavItem(Screen.Notifications, "通知", Icons.Filled.Notifications, Icons.Outlined.Notifications),
    BottomNavItem(Screen.Profile, "我的", Icons.Filled.Person, Icons.Outlined.Person),
)

/**
 * 主导航组件
 */
@Composable
fun AppNavigation() {
    val navController = rememberNavController()
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentDestination = navBackStackEntry?.destination

    // 判断是否显示底部导航栏
    val showBottomBar = currentDestination?.route in bottomNavItems.map { it.screen.route }

    Scaffold(
        bottomBar = {
            if (showBottomBar) {
                NavigationBar {
                    bottomNavItems.forEach { item ->
                        val selected = currentDestination?.hierarchy?.any { it.route == item.screen.route } == true
                        NavigationBarItem(
                            icon = {
                                Icon(
                                    imageVector = if (selected) item.selectedIcon else item.unselectedIcon,
                                    contentDescription = item.label
                                )
                            },
                            label = { Text(item.label) },
                            selected = selected,
                            onClick = {
                                navController.navigate(item.screen.route) {
                                    popUpTo(navController.graph.findStartDestination().id) {
                                        saveState = true
                                    }
                                    launchSingleTop = true
                                    restoreState = true
                                }
                            }
                        )
                    }
                }
            }
        },
        floatingActionButton = {
            if (showBottomBar) {
                FloatingActionButton(
                    onClick = { navController.navigate(Screen.CreatePost.route) },
                    containerColor = MaterialTheme.colorScheme.primary
                ) {
                    Icon(Icons.Default.Add, contentDescription = "发布")
                }
            }
        }
    ) { innerPadding ->
        NavHost(
            navController = navController,
            startDestination = Screen.Feed.route,
            modifier = Modifier.padding(innerPadding)
        ) {
            // 认证页面
            composable(Screen.Login.route) {
                LoginScreen(
                    onNavigateToRegister = { navController.navigate(Screen.Register.route) },
                    onLoginSuccess = {
                        navController.navigate(Screen.Feed.route) {
                            popUpTo(Screen.Login.route) { inclusive = true }
                        }
                    }
                )
            }
            composable(Screen.Register.route) {
                RegisterScreen(
                    onNavigateToLogin = { navController.navigate(Screen.Login.route) },
                    onRegisterSuccess = {
                        navController.navigate(Screen.Feed.route) {
                            popUpTo(Screen.Register.route) { inclusive = true }
                        }
                    }
                )
            }

            // 主页面
            composable(Screen.Feed.route) {
                FeedScreen(
                    onNavigateToPost = { postId ->
                        navController.navigate(Screen.PostDetail.createRoute(postId))
                    },
                    onNavigateToLogin = { navController.navigate(Screen.Login.route) }
                )
            }
            
            // 帖子详情
            composable(Screen.PostDetail.route) {
                PostDetailScreen(
                    onBack = { navController.popBackStack() }
                )
            }

            // 搜索
            composable(Screen.Search.route) {
                SearchScreen(
                    onNavigateToPost = { postId ->
                        navController.navigate(Screen.PostDetail.createRoute(postId))
                    },
                    onNavigateToUser = { userId ->
                        navController.navigate(Screen.UserProfile.createRoute(userId))
                    }
                )
            }

            // 通知
            composable(Screen.Notifications.route) {
                NotificationScreen(
                    onNavigateToPost = { postId ->
                        navController.navigate(Screen.PostDetail.createRoute(postId))
                    },
                    onNavigateToUser = { userId ->
                        navController.navigate(Screen.UserProfile.createRoute(userId))
                    }
                )
            }

            // 我的主页
            composable(Screen.Profile.route) {
                ProfileScreen(
                    onNavigateToLogin = { navController.navigate(Screen.Login.route) },
                    onNavigateToSettings = { /* TODO */ },
                    onNavigateToPost = { postId ->
                        navController.navigate(Screen.PostDetail.createRoute(postId))
                    },
                    onNavigateToUser = { userId ->
                        navController.navigate(Screen.UserProfile.createRoute(userId))
                    }
                )
            }

            // 他人主页
            composable(
                route = Screen.UserProfile.route,
                arguments = listOf(navArgument("userId") { type = NavType.LongType })
            ) {
                ProfileScreen(
                    onNavigateToLogin = { navController.navigate(Screen.Login.route) },
                    onNavigateToSettings = {},
                    onNavigateToPost = { postId ->
                        navController.navigate(Screen.PostDetail.createRoute(postId))
                    },
                    onNavigateToUser = { userId ->
                        navController.navigate(Screen.UserProfile.createRoute(userId))
                    },
                    onBack = { navController.popBackStack() },
                    showBackButton = true
                )
            }

            // 发帖
            composable(Screen.CreatePost.route) {
                CreatePostScreen(
                    onNavigateBack = { navController.popBackStack() },
                    onPostCreated = { navController.popBackStack() }
                )
            }
        }
    }
}
