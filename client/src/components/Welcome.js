import React, { useContext } from 'react'
import { FaConnectdevelop, FaDownload, FaFileDownload } from 'react-icons/fa'
import { Link } from 'react-router-dom'
import { AppContext } from '../context/AppContext'
import { shortenAddress } from '../utils/ShortAdress'

const Welcome = () => {
  const { connectWallet } = useContext(AppContext)

  return (
    <div className="welcome flex max-w-[1500px] m-auto justify-center items-center p-[100px] ">
      <div className="welcome-div-text md:flex-[0.5] flex justify-center px-20 flex-wrap items-center self-center">
        <div className=" w-full text-white text-center text-4xl py-3 font-bold">
          Love Golang
        </div>
        <div className="welcome-button flex items-center cursor-pointer">
        <Link to="/pairs">
          <button
            type="button"
            onClick={connectWallet}
            className="flex flex-row justify-center items-center my-5 bg-[#2952e3] p-3 rounded-lg cursor-pointer hover:bg-[#2546bd]"
          >
            <FaConnectdevelop fontSize={25} className="text-white mr-0" />
            <p className="text-white text-base font-semibold">See pairs</p>
          </button>
          </Link>
        </div>
        <div className="welcome-button flex items-center cursor-pointer mr-5 px-5">
          <Link to="/archives">
            <button
              type="button"
              className="flex flex-row justify-center items-center my-5 bg-[#2952e3] p-3 rounded-lg cursor-pointer hover:bg-[#2546bd]"
            >
              <FaFileDownload fontSize={25} className="text-white mr-0" />
              <p className="text-white text-base font-semibold">
                Download archives
              </p>
            </button>
          </Link>
        </div>
      </div>
    </div>
  )
}

export default Welcome
