import modifyRule from '../assets/modify_rule.png'
import { getAllRules, rule } from '../api/rules';
import { useEffect, useRef, useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';
import { customToastStyle } from '../utils/toast_styles';

interface Props {
    openAddOrUpdateRuleDialog: (rule: rule | null) => void
}

const RulesTable: React.FC<Props> = ({ openAddOrUpdateRuleDialog }) => {
    const [data, setData] = useState<rule[]>();
    const errorShown = useRef(false)

    const fetchRules = async () => {
        try {
            const rules = await getAllRules();
            setData(rules);
            errorShown.current = false
        } catch (error) {
            console.error("Failed to fetch rules:", error);
            if (errorShown.current === false) {
                toast.error("Error: " + error, {
                    style: customToastStyle
                })
                errorShown.current = true
            }
        }
    };

    useEffect(() => {
        fetchRules()
    }, [])



    return (
        <div className="px-8 py-8">
            <table className="table-auto w-full text-left">
                <thead>
                    <tr>
                        <th className="text-left" style={{ width: "50%" }}>Endpoint</th>  {/* Left aligned */}
                        <th className="text-center" style={{ width: "15%" }}>Method</th>  {/* Center aligned */}
                        <th className="text-center" style={{ width: "15%" }}>Strategy</th>  {/* Center aligned */}
                        <th className="text-center" style={{ width: "20%" }}>Modify Rules</th>  {/* Center aligned */}
                    </tr>
                </thead>
                <tbody>
                    {data === undefined ? <div></div> : data.map((item, index) => (
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