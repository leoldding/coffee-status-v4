import React from "react";

const NotFound: React.FC = () => {
    return (
        <>
            <div className="h-dvh flex items-center justify-center gap-2 text-4xl">
                <div className="border-solid border-r-2 pr-2">404</div>
                <div>This page could not be found.</div>
            </div>
        </>
    );
};

export default NotFound;
