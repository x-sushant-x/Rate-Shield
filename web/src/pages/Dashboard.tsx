import { useState } from "react"
import ContentArea from "../components/ContentArea"
import Sidebar from "../components/SideBar"

const Dashboard: React.FC = () => {
    const [selectedPage, setSelectedPage] = useState('API_CONFIGURATION')

    return (
        <div className="flex flex-row bg-global-bg p-3 h-screen overflow-hidden">
            <Sidebar onSelectPage={setSelectedPage} />
            <ContentArea selectedPage={selectedPage} />
        </div>
    )
}

export default Dashboard