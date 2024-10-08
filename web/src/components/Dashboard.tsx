import { useState } from "react"
import ContentArea from "./ContentArea"
import Sidebar from "./SideBar"
import Footer from "./Footer"


const Dashboard: React.FC = () => {
    const [selectedPage, setSelectedPage] = useState('API_CONFIGURATION')

    return (
   <div className="overflow-x-hidden flex flex-col h-screen w-screen lg:flex-row">
      <div className="fixed top-0 -z-20 h-screen w-screen bg-black"></div>

  <div className="lg:grid lg:grid-cols-[20%_80%] w-full h-full">
    <div className="hidden lg:block">
       <Sidebar onSelectPage={setSelectedPage} />           
    </div>
    <div className="w-full h-auto">
      <ContentArea selectedPage={selectedPage} />
    </div>
  </div>
  
  <div className="block lg:hidden">
    <Footer onSelectPage={setSelectedPage}/>
  </div>
   </div>
    )
}

export default Dashboard