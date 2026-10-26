package com.shiyemo.app.ui.post

import android.content.Context
import android.net.Uri
import com.shiyemo.app.BuildConfig
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.shiyemo.app.data.model.dto.CreatePostRequest
import com.shiyemo.app.data.remote.PostApi
import dagger.hilt.android.lifecycle.HiltViewModel
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.MultipartBody
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import javax.inject.Inject

/**
 * 发帖页面 UI 状态
 */
data class CreatePostUiState(
    val type: String = "emotion",
    val title: String = "",
    val content: String = "",
    val tags: List<String> = emptyList(),
    val tagInput: String = "",
    val images: List<Uri> = emptyList(),
    val isAnonymous: Boolean = false,
    val isSubmitting: Boolean = false,
    val error: String? = null,
    val isSuccess: Boolean = false
)

/**
 * 发帖 ViewModel
 */
@HiltViewModel
class CreatePostViewModel @Inject constructor(
    private val postApi: PostApi,
    private val okHttpClient: OkHttpClient,
    @ApplicationContext private val context: Context
) : ViewModel() {

    private val _uiState = MutableStateFlow(CreatePostUiState())
    val uiState: StateFlow<CreatePostUiState> = _uiState.asStateFlow()

    /**
     * 更新帖子类型
     */
    fun updateType(type: String) {
        _uiState.value = _uiState.value.copy(type = type)
    }

    /**
     * 更新标题
     */
    fun updateTitle(title: String) {
        _uiState.value = _uiState.value.copy(title = title)
    }

    /**
     * 更新内容
     */
    fun updateContent(content: String) {
        _uiState.value = _uiState.value.copy(content = content)
    }

    /**
     * 更新标签输入
     */
    fun updateTagInput(input: String) {
        _uiState.value = _uiState.value.copy(tagInput = input)
    }

    /**
     * 添加标签
     */
    fun addTag() {
        val tag = _uiState.value.tagInput.trim()
        if (tag.isNotBlank() && !_uiState.value.tags.contains(tag)) {
            _uiState.value = _uiState.value.copy(
                tags = _uiState.value.tags + tag,
                tagInput = ""
            )
        }
    }

    /**
     * 移除标签
     */
    fun removeTag(tag: String) {
        _uiState.value = _uiState.value.copy(
            tags = _uiState.value.tags - tag
        )
    }

    /**
     * 添加图片
     */
    fun addImage(uri: Uri) {
        if (_uiState.value.images.size < MAX_IMAGES) {
            _uiState.value = _uiState.value.copy(
                images = _uiState.value.images + uri
            )
        }
    }

    /**
     * 移除图片
     */
    fun removeImage(uri: Uri) {
        _uiState.value = _uiState.value.copy(
            images = _uiState.value.images - uri
        )
    }

    /**
     * 切换匿名发布
     */
    fun toggleAnonymous() {
        _uiState.value = _uiState.value.copy(
            isAnonymous = !_uiState.value.isAnonymous
        )
    }

    /**
     * 提交帖子
     */
    fun submitPost() {
        val state = _uiState.value

        // 验证
        if (state.content.isBlank()) {
            _uiState.value = state.copy(error = "请输入帖子内容")
            return
        }

        viewModelScope.launch {
            _uiState.value = state.copy(isSubmitting = true, error = null)

            try {
                // 上传图片（如果有）
                val imageUrls = mutableListOf<String>()
                state.images.forEach { uri ->
                    val url = uploadImage(uri)
                    if (url != null) {
                        imageUrls.add(url)
                    }
                }

                // 创建帖子
                val request = CreatePostRequest(
                    type = state.type,
                    title = state.title,
                    content = state.content,
                    tags = state.tags,
                    images = imageUrls,
                    isAnonymous = state.isAnonymous
                )

                val response = postApi.createPost(request)

                if (response.code == 0) {
                    _uiState.value = _uiState.value.copy(
                        isSubmitting = false,
                        isSuccess = true
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isSubmitting = false,
                        error = response.message
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isSubmitting = false,
                    error = e.message ?: "发布失败"
                )
            }
        }
    }

    /**
     * 上传图片
     */
    private suspend fun uploadImage(uri: Uri): String? {
        return try {
            val inputStream = context.contentResolver.openInputStream(uri) ?: return null
            val bytes = inputStream.readBytes()
            inputStream.close()

            val requestBody = MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart(
                    "file",
                    "image.jpg",
                    bytes.toRequestBody("image/jpeg".toMediaTypeOrNull())
                )
                .build()

            val request = Request.Builder()
                .url("${BuildConfig.API_BASE_URL}upload")
                .post(requestBody)
                .build()

            val response = okHttpClient.newCall(request).execute()
            if (response.isSuccessful) {
                // 解析响应获取图片 URL
                val responseBody = response.body?.string()
                // TODO: 解析响应获取实际的图片 URL
                responseBody?.let {
                    // 临时返回，实际需要解析 JSON
                    java.net.URI(BuildConfig.API_BASE_URL).resolve("/uploads/image.jpg").toString()
                }
            } else {
                null
            }
        } catch (e: Exception) {
            null
        }
    }

    /**
     * 清除错误信息
     */
    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }

    companion object {
        private const val MAX_IMAGES = 9
    }
}
