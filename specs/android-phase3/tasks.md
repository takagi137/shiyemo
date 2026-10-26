# tasks：Android App Phase 3 — 社交功能

## Phase 1：数据层

- [x] 任务 1.1：新增 SearchApi 接口
  - 文件：`data/remote/SearchApi.kt`
  - 验收：`GET /search?q=` 定义正确
  - Spec：01-search-success

- [x] 任务 1.2：新增 NotificationApi 接口
  - 文件：`data/remote/NotificationApi.kt`
  - 验收：列表/已读/全部已读三个接口定义正确
  - Spec：03-notification-list

- [x] 任务 1.3：扩展 UserApi（关注相关）
  - 文件：`data/remote/UserApi.kt`
  - 验收：follow/unfollowing/followers 四个接口定义正确
  - Spec：04-follow-unfollow, 05-follow-list

- [x] 任务 1.4：新增 DTO（SearchResultDto, NotificationDto）
  - 文件：`data/model/dto/PostDto.kt`, `UserDto.kt`
  - 验收：字段与后端 API 响应对齐
  - Spec：01-search-success, 03-notification-list

- [x] 任务 1.5：新增 Repository（Search, Notification）
  - 文件：`data/repository/SearchRepository.kt`, `NotificationRepository.kt`
  - 验收：方法签名与 API 接口一致
  - Spec：01-search-success, 03-notification-list

- [x] 任务 1.6：扩展 ProfileRepository（关注相关）
  - 文件：`data/repository/ProfileRepository.kt`
  - 验收：新增 followUser/unfollowUser/getFollowing/getFollowers
  - Spec：04-follow-unfollow, 05-follow-list

- [x] 任务 1.7：注册新 API 到 DI
  - 文件：`di/NetworkModule.kt`
  - 验收：SearchApi、NotificationApi 通过 Hilt 注入
  - Spec：全部

## Phase 2：ViewModel 层

- [x] 任务 2.1：SearchViewModel
  - 文件：`ui/search/SearchViewModel.kt`
  - 验收：防抖 500ms、分页加载、清空搜索
  - Spec：01-search-success, 02-search-debounce

- [x] 任务 2.2：NotificationViewModel
  - 文件：`ui/notification/NotificationViewModel.kt`
  - 验收：列表加载、标记已读、全部已读
  - Spec：03-notification-list

- [x] 任务 2.3：ProfileViewModel 增强
  - 文件：`ui/profile/ProfileViewModel.kt`
  - 验收：支持他人主页、关注/取关、关注/粉丝列表
  - Spec：04-follow-unfollow, 05-follow-list

## Phase 3：UI 层

- [x] 任务 3.1：SearchScreen 完整实现
  - 文件：`ui/search/SearchScreen.kt`
  - 验收：搜索栏 + 结果卡片 + 无限滚动 + 空状态
  - Spec：01-search-success, 02-search-debounce

- [x] 任务 3.2：NotificationScreen 完整实现
  - 文件：`ui/notification/NotificationScreen.kt`
  - 验收：通知列表 + 未读标记 + 全部已读 + 类型图标
  - Spec：03-notification-list

- [x] 任务 3.3：ProfileScreen 增强
  - 文件：`ui/profile/ProfileScreen.kt`
  - 验收：关注按钮 + 用户列表 Tab + 他人主页返回按钮
  - Spec：04-follow-unfollow, 05-follow-list

- [x] 任务 3.4：导航更新
  - 文件：`ui/navigation/AppNavigation.kt`
  - 验收：新增 user/{userId} 路由，搜索/通知导航回调正确
  - Spec：04-follow-unfollow, 05-follow-list

## Phase 4：单测

- [x] 任务 4.1：SearchRepository 单测
  - 文件：`test/.../SearchRepositoryTest.kt`
  - 验收：搜索成功、空结果、类型过滤、分页、服务端错误、网络异常（6 个用例）
  - Spec：01-search-success

- [x] 任务 4.2：NotificationRepository 单测
  - 文件：`test/.../NotificationRepositoryTest.kt`
  - 验收：列表加载、空列表、分页、发送者信息、标记已读、全部已读、服务端错误、网络异常（8 个用例）
  - Spec：03-notification-list

- [x] 任务 4.3：ProfileRepository 关注功能单测
  - 文件：`test/.../ProfileRepositoryTest.kt`
  - 验收：关注成功/重复、取关成功/未关注、关注列表/分页/空、粉丝列表/分页/空、用户信息、收藏、网络异常（14 个用例）
  - Spec：04-follow-unfollow, 05-follow-list
