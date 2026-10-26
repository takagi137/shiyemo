package com.shiyemo.app.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.local.prefs.TokenManager
import com.shiyemo.app.data.model.dto.LoginEmailRequest
import com.shiyemo.app.data.model.dto.RegisterEmailRequest
import com.shiyemo.app.data.remote.AuthApi
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

/**
 * 认证页面 UI 状态
 */
data class AuthUiState(
    val isLoading: Boolean = false,
    val error: String? = null,
    val isLoggedIn: Boolean = false
)

/**
 * 认证 ViewModel
 *
 * 处理登录、注册逻辑，连接后端 API
 */
@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authApi: AuthApi,
    private val tokenManager: TokenManager
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    /**
     * 邮箱登录
     */
    fun login(email: String, password: String) {
        if (email.isBlank() || password.isBlank()) {
            _uiState.value = _uiState.value.copy(error = "请输入邮箱和密码")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val response = authApi.loginEmail(
                    LoginEmailRequest(email = email, password = password)
                )

                if (response.isSuccess && response.data != null) {
                    val authData = response.data!!
                    tokenManager.saveAuthInfo(
                        token = authData.token,
                        userId = authData.user.id,
                        nickname = authData.user.nickname,
                        avatar = authData.user.avatar
                    )
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        isLoggedIn = true
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = response.message.ifBlank { "登录失败" }
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "网络错误，请重试"
                )
            }
        }
    }

    /**
     * 邮箱注册
     */
    fun register(email: String, password: String, nickname: String) {
        if (email.isBlank() || password.isBlank() || nickname.isBlank()) {
            _uiState.value = _uiState.value.copy(error = "请填写完整信息")
            return
        }

        if (password.length < 6) {
            _uiState.value = _uiState.value.copy(error = "密码至少6位")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val response = authApi.registerEmail(
                    RegisterEmailRequest(email = email, password = password, nickname = nickname)
                )

                if (response.isSuccess && response.data != null) {
                    val authData = response.data!!
                    tokenManager.saveAuthInfo(
                        token = authData.token,
                        userId = authData.user.id,
                        nickname = authData.user.nickname,
                        avatar = authData.user.avatar
                    )
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        isLoggedIn = true
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = response.message.ifBlank { "注册失败" }
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "网络错误，请重试"
                )
            }
        }
    }

    /**
     * 清除错误信息
     */
    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}
