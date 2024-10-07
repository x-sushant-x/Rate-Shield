import infoIcon from './../assets/Info Squared.svg'
import githubIcon from './../assets/GitHub.svg'
import twitterIcon from './../assets/Twitter.svg'
import linkedinIcon from './../assets/LinkedIn.svg'


interface SidebarProps {
    onSelectPage: (page: string) => void
}

const Footer: React.FC<SidebarProps> = ({ onSelectPage })=> {
  return (
       <footer className="w-full flex-row items-center justify-center bg-sidebar-bg text-gray-300">
            <div className='flex flex-row justify-between mx-10 pt-7 items-center'>
            <div className='flex items-center cursor-pointer'
                    onClick={() => onSelectPage('ABOUT')}>
                    <img src={infoIcon} />
                </div>

                <div className='flex items-center  cursor-pointer'
                    onClick={() => onSelectPage('TWITTER')}>
                    <img src={twitterIcon} />
                </div>

                <div className='flex items-center cursor-pointer'
                    onClick={() => onSelectPage('LINKEDIN')}>
                    <img src={linkedinIcon} />
                </div>

                <div className='flex items-center cursor-pointer'
                    onClick={() => onSelectPage('GITHUB')}>
                    <img src={githubIcon} />
                </div>
            </div>
            <div className='flex items-center p-4 justify-center text-sm'>
                Created for fun by  <b>  Sushant</b>
            </div>
        </footer>
  )
}

export default Footer