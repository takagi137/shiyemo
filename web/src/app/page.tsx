import Link from 'next/link'
import DailyScene from '@/components/DailyScene'

export default function HomePage() {
  return (
    <div className="min-h-screen">
      {/* 全屏风景图 + 温暖文案 */}
      <DailyScene />

      {/* 产品理念 */}
      <section className="py-20 bg-white">
        <div className="max-w-4xl mx-auto px-6 text-center">
          <h2 className="text-3xl md:text-4xl font-light text-gray-800 mb-8">
            在这里，没有人在催你
          </h2>
          <p className="text-lg text-gray-500 leading-relaxed mb-12">
            慢慢来，你已经在路上了。
          </p>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-12">
            <div className="p-6">
              <div className="text-4xl mb-4">🌊</div>
              <h3 className="text-lg font-medium text-gray-800 mb-2">不评判，只记录</h3>
              <p className="text-gray-500 text-sm">
                你的情绪、你的思考、你的尝试，都值得被看见。
              </p>
            </div>
            <div className="p-6">
              <div className="text-4xl mb-4">🤝</div>
              <h3 className="text-lg font-medium text-gray-800 mb-2">你不是一个人</h3>
              <p className="text-gray-500 text-sm">
                很多人和你一样，正在经历同样的困惑和不安。
              </p>
            </div>
            <div className="p-6">
              <div className="text-4xl mb-4">🌱</div>
              <h3 className="text-lg font-medium text-gray-800 mb-2">成长需要时间</h3>
              <p className="text-gray-500 text-sm">
                学会做饭、学会独处、学会重新开始，都是成果。
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-16 bg-gradient-to-r from-blue-50 to-purple-50">
        <div className="max-w-4xl mx-auto px-6 text-center">
          <p className="text-xl text-gray-600 mb-8">
            准备好了吗？
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-4">
            <Link
              href="/register"
              className="px-8 py-3 bg-blue-600 text-white font-medium rounded-full hover:bg-blue-700 transition-colors"
            >
              加入社区
            </Link>
            <Link
              href="/feed"
              className="px-8 py-3 bg-white text-gray-700 font-medium rounded-full hover:bg-gray-50 transition-colors border border-gray-200"
            >
              浏览动态
            </Link>
          </div>
        </div>
      </section>

      {/* 页脚 */}
      <footer className="py-8 bg-gray-50 border-t border-gray-100">
        <div className="max-w-4xl mx-auto px-6 text-center">
          <p className="text-gray-400 text-sm">
            © 2024 失业么 — AI 时代个体的能力重建与收入重建社区
          </p>
        </div>
      </footer>
    </div>
  )
}
