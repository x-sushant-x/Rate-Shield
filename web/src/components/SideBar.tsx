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
        <div className="w-64 h-screen bg-sidebar-bg text-gray-300 rounded-xl ">
            <div className="p-6">
                <img src={logo} />
            </div>
            <ul>
                <div className='flex items-center ml-7 mt-4 cursor-pointer'
                    onClick={() => onSelectPage('API_CONFIGURATION')}>
                    <img src={apiIcon} />
                    <li className="ml-4 hover:bg-gray-700 text-sm">API Configuration</li>
                </div>

                <div className='flex items-center ml-7 mt-4 cursor-pointer'
                    onClick={() => onSelectPage('ABOUT')}>
                    <img src={infoIcon} />
                    <li className="ml-4 hover:bg-gray-700 text-sm">About</li>
                </div>

                <div className='flex items-center ml-7 mt-4 cursor-pointer'
                    onClick={() => onSelectPage('TWITTER')}>
                    <img src={twitterIcon} />
                    <li className="ml-4 hover:bg-gray-700 text-sm">Follow on X</li>
                </div>

                <div className='flex items-center ml-7 mt-4 cursor-pointer'
                    onClick={() => onSelectPage('LINKEDIN')}>
                    <img src={linkedinIcon} />
                    <li className="ml-4 hover:bg-gray-700 text-sm">Follow on LinkedIn</li>
                </div>

                <div className='flex items-center ml-7 mt-4 cursor-pointer'
                    onClick={() => onSelectPage('GITHUB')}>
                    <img src={githubIcon} />
                    <li className="ml-4 hover:bg-gray-700 text-sm">Follow on GitHub</li>
                </div>
            </ul>

            <div className='text-sm ml-6 mt-[40rem]'>
                Created for fun by <b>Sushant</b>
            </div>
        </div>
    );
};

export default Sidebar
