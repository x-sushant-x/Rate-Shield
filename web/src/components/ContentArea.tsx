import APIConfiguration from "../pages/APIConfiguration";

interface ContentAreaProps {
    selectedPage: string
}

const ContentArea: React.FC<ContentAreaProps> = ({ selectedPage }) => {
    return (
        <div className="flex-1">
            {selectedPage === 'API_CONFIGURATION' && <APIConfiguration />}
            {selectedPage === 'ABOUT' && <div>About Content</div>}
            {selectedPage === 'TWITTER' && <div>Tweet on X Content</div>}
            {selectedPage === 'LINKEDIN' && <div>Follow on LinkedIn Content</div>}
            {selectedPage === 'GITHUB' && <div>Follow on GitHub Content</div>}
        </div>
    );
};

export default ContentArea;
