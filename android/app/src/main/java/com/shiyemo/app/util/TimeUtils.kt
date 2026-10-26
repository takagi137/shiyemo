package com.shiyemo.app.util

import java.text.SimpleDateFormat
import java.util.Calendar
import java.util.Date
import java.util.Locale
import java.util.concurrent.TimeUnit

/**
 * 时间格式化工具
 *
 * 将 ISO 8601 格式的时间字符串转换为友好的相对时间显示
 */
object TimeUtils {

    private val isoFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss", Locale.getDefault())
    private val dateTimeFormat = SimpleDateFormat("yyyy-MM-dd HH:mm", Locale.getDefault())
    private val dateFormat = SimpleDateFormat("yyyy-MM-dd", Locale.getDefault())

    /**
     * 格式化为相对时间
     *
     * - 1分钟内：刚刚
     * - 1小时内：X分钟前
     * - 24小时内：X小时前
     * - 7天内：X天前
     * - 同年：MM-dd HH:mm
     * - 跨年：yyyy-MM-dd
     */
    fun formatRelative(timeStr: String): String {
        if (timeStr.isBlank()) return ""

        val date = parseDate(timeStr) ?: return timeStr.take(10)
        val now = Date()
        val diff = now.time - date.time

        // 未来时间（可能时区问题）直接显示日期
        if (diff < 0) return dateTimeFormat.format(date)

        val minutes = TimeUnit.MILLISECONDS.toMinutes(diff)
        val hours = TimeUnit.MILLISECONDS.toHours(diff)
        val days = TimeUnit.MILLISECONDS.toDays(diff)

        return when {
            minutes < 1 -> "刚刚"
            minutes < 60 -> "${minutes}分钟前"
            hours < 24 -> "${hours}小时前"
            days < 7 -> "${days}天前"
            else -> {
                val calNow = Calendar.getInstance()
                val calDate = Calendar.getInstance().apply { time = date }
                if (calNow.get(Calendar.YEAR) == calDate.get(Calendar.YEAR)) {
                    dateTimeFormat.format(date)
                } else {
                    dateFormat.format(date)
                }
            }
        }
    }

    /**
     * 格式化为完整日期时间
     */
    fun formatFull(timeStr: String): String {
        if (timeStr.isBlank()) return ""
        val date = parseDate(timeStr) ?: return timeStr
        return dateTimeFormat.format(date)
    }

    /**
     * 解析时间字符串，支持多种格式
     */
    private fun parseDate(timeStr: String): Date? {
        // 尝试 ISO 8601 格式
        isoFormat.parse(timeStr)?.let { return it }
        // 尝试 yyyy-MM-dd HH:mm 格式
        dateTimeFormat.parse(timeStr)?.let { return it }
        // 尝试 yyyy-MM-dd 格式
        dateFormat.parse(timeStr)?.let { return it }
        // 尝试带时区的 ISO 格式（截取前19位）
        if (timeStr.length > 19) {
            isoFormat.parse(timeStr.substring(0, 19))?.let { return it }
        }
        return null
    }
}
