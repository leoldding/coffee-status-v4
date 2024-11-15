import React from "react";
import "./404.css";

const NotFound: React.FC = () => {
    return (
        <>
            <div className="notfound-container">
                <div className="notfound-border">404</div>
                <div>This page could not be found.</div>
            </div>
        </>
    );
};

export default NotFound;
