import React from "react";
import { BrowserRouter as Router, Route, Routes  } from "react-router-dom";
import { Home, Admin, NotFound } from "./pages";

const App: React.FC = () => {
    return (
        <>
            <Router basename="/coffee">
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/admin" element={<Admin />} />
                    <Route path="/*" element={<NotFound />} />
                </Routes>
            </Router>
        </>
    )
}

export default App;
