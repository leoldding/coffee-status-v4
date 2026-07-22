import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { FaHome } from "react-icons/fa";
import { MdCheck, MdClose } from "react-icons/md";
import { api } from "../api/client";

const Admin: React.FC = () => {
    const [authenticated, setAuthenticated] = useState<boolean>(false);
    const [updateStatus, setUpdateStatus] = useState<number>(0); // 0=hidden, 1=loading, 2=success, 3=failure

    useEffect(() => {
        const check = async () => {
            try {
                await api("/api/v1/auth/check");
                setAuthenticated(true);
            } catch {
                // api() already redirected.
            }
        };

        check();
    }, []);

    const putStatus = async (status: string) => {
        setUpdateStatus(1);
        try {
            const response = await api("/api/v1/coffee/status", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    value: status,
                }),
            });

            if (!response.ok) {
                throw new Error("Error updating status");
            }
            setUpdateStatus(2);
        } catch (error) {
            setUpdateStatus(3);
            console.error(error);
        }
    };

    const updateStatusElement = () => {
        if (updateStatus === 0) {
            return
        }
        if (updateStatus === 1) {
            return <span className="loading loading-spinner"></span>
        }
        if (updateStatus === 2) {
            return <span className="btn btn-success btn-circle btn-xs" onClick={() => setUpdateStatus(0)}><MdCheck /></span>
        }
        if (updateStatus === 3) {
            return <span className="btn btn-error btn-circle btn-xs" onClick={() => setUpdateStatus(0)}><MdClose /></span>
        }
    }

    return (
        authenticated &&
        <>
            <div className="h-dvh flex flex-col items-center justify-center gap-8">
                <button className="btn btn-xl w-48 btn-success" onClick={(() => putStatus("yes"))}>Yes</button>
                <button className="btn btn-xl w-48 btn-warning" onClick={(() => putStatus("otw"))}>On the way</button>
                <button className="btn btn-xl w-48 btn-error" onClick={(() => putStatus("no"))}>No</button>
                {updateStatusElement()}
            </div>
            <div className="fixed bottom-6 right-6">
                <Link to="/">
                    <button className="flex items-center justify-center btn btn-lg btn-soft btn-accent btn-circle">
                        <FaHome />
                    </button>
                </Link>
            </div>
        </>
    )
}

export default Admin;
