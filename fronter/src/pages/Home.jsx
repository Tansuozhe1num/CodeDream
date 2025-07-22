import React from 'react'
import { motion } from 'framer-motion'
import { ReactComponent as CodeIcon } from '../assets/code-icon.svg'
import { ReactComponent as TrophyIcon } from '../assets/trophy-icon.svg'
import { ReactComponent as NoteIcon } from '../assets/note-icon.svg'

export default function Home() {
    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex flex-col">
            {/* 导航栏 */}
            <nav className="bg-white shadow-md">
                <div className="container mx-auto px-6 py-4 flex justify-between items-center">
                    <div className="text-2xl font-bold text-indigo-700">Codedream</div>
                    <div className="space-x-4">
                        <a href="/" className="text-gray-700 hover:text-indigo-600">首页</a>
                        <a href="/problems" className="text-gray-700 hover:text-indigo-600">题目</a>
                        <a href="/ranking" className="text-gray-700 hover:text-indigo-600">排行榜</a>
                        <a href="/notes" className="text-gray-700 hover:text-indigo-600">笔记</a>
                        <a href="/login" className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition">登录</a>
                    </div>
                </div>
            </nav>

            {/* Hero 区块 */}
            <motion.section
                className="flex-1 flex flex-col justify-center items-center text-center px-4"
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.8 }}
            >
                <h1 className="text-5xl md:text-6xl font-extrabold text-indigo-800 mb-4">
                    高效在线刷题平台
                </h1>
                <p className="text-lg md:text-xl text-gray-600 max-w-2xl mb-8">
                    Codedream 聚合丰富题库、实时排行和笔记分享，助你挑战算法巅峰。
                </p>
                <div className="flex space-x-6">
                    <a
                        href="/problems"
                        className="px-8 py-3 bg-indigo-600 text-white rounded-lg shadow hover:bg-indigo-700 transition"
                    >
                        开始刷题
                    </a>
                    <a
                        href="/ranking"
                        className="px-8 py-3 border border-indigo-600 text-indigo-600 rounded-lg hover:bg-indigo-50 transition"
                    >
                        查看排行
                    </a>
                </div>
            </motion.section>

            {/* 特色功能区 */}
            <section className="py-16 bg-white">
                <div className="container mx-auto grid grid-cols-1 md:grid-cols-3 gap-8 px-6">
                    <motion.div
                        className="p-6 bg-indigo-50 rounded-2xl shadow-lg text-center"
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.2 }}
                    >
                        <CodeIcon className="w-12 h-12 mx-auto mb-4 text-indigo-600" />
                        <h3 className="text-xl font-semibold mb-2">海量题库</h3>
                        <p className="text-gray-600">覆盖算法、数据结构，多种难度，满足不同阶段需求。</p>
                    </motion.div>
                    <motion.div
                        className="p-6 bg-indigo-50 rounded-2xl shadow-lg text-center"
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.4 }}
                    >
                        <TrophyIcon className="w-12 h-12 mx-auto mb-4 text-indigo-600" />
                        <h3 className="text-xl font-semibold mb-2">实时排行</h3>
                        <p className="text-gray-600">全球刷题排行，激发学习动力，与高手同场竞技。</p>
                    </motion.div>
                    <motion.div
                        className="p-6 bg-indigo-50 rounded-2xl shadow-lg text-center"
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.6 }}
                    >
                        <NoteIcon className="w-12 h-12 mx-auto mb-4 text-indigo-600" />
                        <h3 className="text-xl font-semibold mb-2">笔记分享</h3>
                        <p className="text-gray-600">高效记录解题思路，笔记可公开分享或私有保存。</p>
                    </motion.div>
                </div>
            </section>

            {/* 页脚 */}
            <footer className="bg-indigo-800 text-white py-4">
                <div className="container mx-auto text-center text-sm">
                    &copy; {new Date().getFullYear()} Codedream. 保留所有权利。
                </div>
            </footer>
        </div>
    )
}
