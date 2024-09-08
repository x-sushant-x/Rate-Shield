import modifyRule from '../assets/modify_rule.png'
import { rule } from '../api/rules';
import { Toaster } from 'react-hot-toast';

interface Props {
    openAddOrUpdateRuleDialog: (rule: rule | null) => void
    rulesData: rule[] | undefined
}

const RulesTable: React.FC<Props> = ({ openAddOrUpdateRuleDialog, rulesData }) => {
    return (
        <div className="px-8 py-8">
            <table className="table-auto w-full text-left">
                <thead>
                    <tr>
                        <th className="text-left" style={{ width: "63%" }}>Endpoint</th>  {/* Left aligned */}
                        <th className="text-center" style={{ width: "5%" }}>Method</th>  {/* Center aligned */}
                        <th className="text-center" style={{ width: "10%" }}>Strategy</th>  {/* Center aligned */}
                        <th className="text-center" style={{ width: "15%" }}>Modify</th>  {/* Center aligned */}
                    </tr>
                </thead>
                <tbody>
                    {rulesData === undefined ? <div>
                        Unable to fetch rules.
                    </div> : rulesData.map((item, index) => (
                        <>
                            <tr key={index}>
                                <td style={{ width: "50%" }} className='pt-6'>{item.endpoint}</td>
                                <td className="text-center pt-6" style={{ width: "15%" }}>{item.http_method}</td>  {/* Center aligned */}
                                <td className="text-center pt-6" style={{ width: "15%" }}>{item.strategy}</td>  {/* Center aligned */}
                                <td className="text-center pt-6" style={{ width: "20%" }}>
                                    <center>
                                        <img src={modifyRule} className='cursor-pointer' onClick={() => {
                                            openAddOrUpdateRuleDialog(item)
                                        }} />
                                    </center>
                                </td>
                            </tr>

                            <tr>
                                <td colSpan={4}>
                                    <hr className="border-gray-300 mt-4" />
                                </td>
                            </tr>
                        </>
                    ))}
                </tbody>
            </table>
            <Toaster position="bottom-right" reverseOrder={false} />
        </div>
    )
}

export default RulesTable