import React, { useContext, useEffect } from 'react'
import datas from '../utils/datas'
import { Loader } from '../components/'
import { AppContext } from '../context/AppContext'

const companyCommonStyles =
  'min-h-[70px] sm:px-0 px-2 sm:min-w-[220px] flex justify-center items-center border-[0.5px] border-gray-400 text-sm font-light text-white'
const Pairs = () => {
  //call getData function every minute
  const { serverData } = useContext(AppContext)
  console.log('serverData =>', serverData.length)
  let dataLength = serverData.length
  console.log('dataLength =>', dataLength)
  return (
    <div className="flex w-full justify-center items-center">
      <div className="flex mf:flex-row flex-col items-center justify-between  py-12 px-4">
        <div className="flex flex-1 justify-center items-center flex-col  mf:ml-30">
          <h1 className="text-3xl text-white  py-1">Pairs </h1>
          <div className="grid sm:grid-cols-8 grid-cols-2 w-full mt-5">
            <div
              className={`rounded-tl-2xl sm:rounded-bl-2xl ${companyCommonStyles}`}
            >
              <p className="font-bold text-[20px]">Pairs</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">High</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">Low</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">VOLUME</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">Ask</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">Bid</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">Price</p>
            </div>
            <div
              className={`sm:rounded-tr-2xl rounded-br-2xl ${companyCommonStyles}`}
            >
              <p className="font-bold text-[20px]">Time</p>
            </div>
          </div>
          {dataLength == 0 ? (
            <>
              <h1 className="text-3xl text-white  py-5">Loading... </h1>
              <Loader />
            </>
          ) : (
            [...serverData[0]].map((pair, i) => (
              <div className="grid sm:grid-cols-8 w-full" key={i}>
                <>
                  <div
                    className={`rounded-tl-2xl sm:rounded-bl-2xl ${companyCommonStyles}`}
                  >
                    {pair.ASSETS.wsname}
                  </div>
                  <div className={companyCommonStyles}>{pair.HIGH}</div>
                  <div className={companyCommonStyles}>{pair.LOW}</div>
                  <div className={companyCommonStyles}>{pair.VOLUME}</div>
                  <div className={companyCommonStyles}>{pair.PRIX_ACHAT}</div>
                  <div className={companyCommonStyles}>{pair.PRIX_VENTE}</div>
                  <div className={companyCommonStyles}>{pair.PRICE}</div>
                  <div
                    className={`sm:rounded-tr-2xl rounded-br-2xl ${companyCommonStyles}`}
                  >
                    {pair.TIMESTAMP}
                  </div>
                </>
              </div>
            ))
          )}
          <h1 className="text-3xl text-white text-gradient py-10">
            Second Tour
          </h1>
        </div>
      </div>
    </div>
  )
}

export default Pairs
