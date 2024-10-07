import { useState } from "react"
import ContentArea from "./ContentArea"
import Sidebar from "./SideBar"
import Footer from "./Footer"


const Dashboard: React.FC = () => {
    const [selectedPage, setSelectedPage] = useState('API_CONFIGURATION')

    return (
        <div className="w-full h-full flex-col lg:flex-row bg-global-bg">
            <div className="hidden lg:block">
            <Sidebar onSelectPage={setSelectedPage} />   
            </div>
            <ContentArea selectedPage={selectedPage} />
            <div className="block lg:hidden">
            <Footer onSelectPage={setSelectedPage}>
            </Footer>
            </div>   
        </div>
    )
}

export default Dashboard