import { useState } from "react"
import ContentArea from "./ContentArea"
import Sidebar from "./SideBar"
import Footer from "./Footer"


const Dashboard: React.FC = () => {
    const [selectedPage, setSelectedPage] = useState('API_CONFIGURATION')

    return (
        <div className="w-screen flex flex-col lg:flex-row bg-global-bg space-x-4 p-3">
            <div className="hidden lg:block">
            <Sidebar onSelectPage={setSelectedPage} />   
            </div>
            <ContentArea selectedPage={selectedPage} />
            <div className="block lg:hidden">
            <Footer>
            </Footer>
            </div>   
        </div>
    )
}

export default Dashboard