import { useEffect } from "react";
import About from "../pages/About";
import APIConfiguration from "../pages/APIConfiguration";

interface ContentAreaProps {
    selectedPage: string
}

const openInNewTab = (url: string) => {
    window.open(url, "_blank", "noreferrer")
}

const ContentArea: React.FC<ContentAreaProps> = ({ selectedPage }) => {
    useEffect(() => {
        if (selectedPage === "TWITTER") {
            openInNewTab('https://x.com/SushantDhiman17')
        }

        if (selectedPage === "LINKEDIN") {
            openInNewTab('https://linkedin.com/in/sushant102004')
        }

        if (selectedPage === "GITHUB") {
            openInNewTab('https://github.com/x-sushant-x')
        }
    }, [selectedPage])


    return (
        <div className="flex-1">
            {selectedPage === 'API_CONFIGURATION' && <APIConfiguration />}
            {selectedPage === 'ABOUT' && <About />}
            {selectedPage === 'TWITTER' && <APIConfiguration />}
            {selectedPage === 'LINKEDIN' && <APIConfiguration />}
            {selectedPage === 'GITHUB' && <APIConfiguration />}
        </div>
    );
};

export default ContentArea;
