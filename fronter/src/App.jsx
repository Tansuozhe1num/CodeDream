import React from 'react';
import { Routes, Route, Link, Outlet } from 'react-router-dom';
import Home from './pages/Home';
import Login from './pages/Login';
import Problems from './pages/Problems';
import Ranking from './pages/Ranking';
import Notes from './pages/Notes';
import Solutions from './pages/Solutions';
import 'animate.css';

export default function App() {
    return (
        <>
            <nav className="navbar">
                <Link to="/">首页</Link>
                <Link to="/problems">题目</Link>
                <Link to="/ranking">排行榜</Link>
                <Link to="/notes">笔记</Link>
                <Link to="/solutions">解题</Link>
                <Link to="/login">登录</Link>
            </nav>
            <div className="content">
                <Outlet />
            </div>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/problems" element={<Problems />} />
                <Route path="/ranking" element={<Ranking />} />
                <Route path="/notes" element={<Notes />} />
                <Route path="/solutions" element={<Solutions />} />
            </Routes>
        </>
    );
}