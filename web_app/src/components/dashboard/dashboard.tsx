import ContentArea from "./content_area"
import Sidebar from "./sidebar"

const Dashboard: React.FC = () => {
    return <>
        <div className="flex flex-row bg-global-bg space-x-4 p-3">
            <Sidebar />
            <ContentArea />
        </div>
    </>
}

export default Dashboard