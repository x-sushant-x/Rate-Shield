import logo from './../assets/logo.svg'
import apiIcon from './../assets/API.svg'
import infoIcon from './../assets/Info Squared.svg'
import githubIcon from './../assets/GitHub.svg'
import twitterIcon from './../assets/Twitter.svg'
import linkedinIcon from './../assets/LinkedIn.svg'

interface SidebarProps {
    onSelectPage: (page: string) => void
}

const Sidebar: React.FC<SidebarProps> = ({ onSelectPage }) => {
    return (
        <div className="w-64 bg-sidebar-bg text-gray-300 rounded-xl flex flex-col justify-between mr-4">
            <div>
                <div className="p-6">
                    <img src={logo} />
                </div>
                <ul>
            {[
                { label: 'API Configuration', icon: apiIcon, page: 'API_CONFIGURATION' },
                { label: 'About', icon: infoIcon, page: 'ABOUT' },
                { label: 'Follow on X', icon: twitterIcon, page: 'TWITTER' },
                { label: 'Follow on LinkedIn', icon: linkedinIcon, page: 'LINKEDIN' },
                { label: 'Follow on GitHub', icon: githubIcon, page: 'GITHUB' }
            ].map((item, index) => (
                <li
                    key={index}
                    className="flex items-center ml-7 mt-4 cursor-pointer text-sm hover:bg-gray-700"
                    onClick={() => onSelectPage(item.page)}
                >
                    <img src={item.icon} alt={`${item.label} Icon`} className="mr-4" />
                    {item.label}
                </li>
            ))}
        </ul>
            </div>

            <div className='text-sm ml-6 mb-4'>
                Still Under Development
            </div>
        </div>
    );
};

export default Sidebar;
