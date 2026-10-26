package com.shiyemo.app.data.local.prefs

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.longPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

// DataStore 实例
private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "shiyemo_prefs")

/**
 * Token 和用户信息本地存储
 */
@Singleton
class TokenManager @Inject constructor(
    @ApplicationContext private val context: Context
) {
    companion object {
        private val TOKEN_KEY = stringPreferencesKey("auth_token")
        private val USER_ID_KEY = longPreferencesKey("user_id")
        private val NICKNAME_KEY = stringPreferencesKey("nickname")
        private val AVATAR_KEY = stringPreferencesKey("avatar")
    }

    /**
     * Token Flow
     */
    val tokenFlow: Flow<String?> = context.dataStore.data.map { prefs ->
        prefs[TOKEN_KEY]
    }

    /**
     * 用户 ID Flow
     */
    val userIdFlow: Flow<Long?> = context.dataStore.data.map { prefs ->
        prefs[USER_ID_KEY]
    }

    /**
     * 保存认证信息
     */
    suspend fun saveAuthInfo(token: String, userId: Long, nickname: String, avatar: String) {
        context.dataStore.edit { prefs ->
            prefs[TOKEN_KEY] = token
            prefs[USER_ID_KEY] = userId
            prefs[NICKNAME_KEY] = nickname
            prefs[AVATAR_KEY] = avatar
        }
    }

    /**
     * 清除认证信息（退出登录）
     */
    suspend fun clearAuthInfo() {
        context.dataStore.edit { prefs ->
            prefs.remove(TOKEN_KEY)
            prefs.remove(USER_ID_KEY)
            prefs.remove(NICKNAME_KEY)
            prefs.remove(AVATAR_KEY)
        }
    }

    /**
     * 获取 Token（同步方法）
     */
    suspend fun getToken(): String? {
        var token: String? = null
        tokenFlow.collect {
            token = it
            return@collect
        }
        return token
    }

    /**
     * 清除 Token（退出登录）
     */
    suspend fun clearToken() {
        clearAuthInfo()
    }

    /**
     * 检查是否已登录
     */
    fun isLoggedIn(): Flow<Boolean> {
        return tokenFlow.map { !it.isNullOrEmpty() }
    }
}
