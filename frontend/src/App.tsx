import React from "react";
import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import { Home, Admin, NotFound } from "./pages";

const App: React.FC = () => {
    return (
        <>
            <Router basename="/coffee">
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/admin" element={<Admin />} />
                    <Route path="/404" element={<NotFound />} />
                    <Route path="*" element={<Navigate to="/404" />} />
                </Routes>
            </Router>
        </>
    )
}

export default App;
