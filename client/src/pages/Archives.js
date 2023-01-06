import React, { useContext, useEffect } from 'react'
import datas from '../utils/datas'
import { Loader } from '../components/'
import { AppContext } from '../context/AppContext'

const companyCommonStyles =
  'min-h-[70px] sm:px-0 px-2 sm:min-w-[220px] flex justify-center items-center border-[0.5px] border-gray-400 text-sm font-light text-white'
const Archives = () => {
  //call getData function every minute
  const { serverFiles } = useContext(AppContext)
  let dataLength = serverFiles.length
  console.log('dataLength =>', dataLength)

  const parts = 'assets_2023-01-05_0H.json'.split('_')
  console.log('parts =>', parts)

  function getTimestamp(filename) {
    // console.log('filename =>', filename)
    // const parts = filename.split('_')
    // const date = parts[1]
    // const time = parts[2]
    // let trueTime = date + '' + time + ':00' + ' Hrs'
    // console.log('trueTime =>', trueTime)
    // return trueTime

    const firstUnderscoreIndex = filename.indexOf('_')

    // Find the position of the period after the timestamp
    const periodIndex = filename.indexOf('.')

    // Get the substring with the timestamp
    const timestamp = filename.substring(firstUnderscoreIndex + 1, periodIndex)
    return timestamp
  }

  getTimestamp('assets_2023-01-05_0H.json')

  return (
    <div className="flex w-full justify-center items-center">
      <div className="flex mf:flex-row flex-col items-center justify-between  py-12 px-4">
        <div className="flex flex-1 justify-center items-center flex-col  mf:ml-30">
          <h1 className="text-3xl text-white  py-1"> Archives </h1>
          <div className="grid sm:grid-cols-3 grid-cols-2 w-full mt-5">
            <div
              className={`rounded-tl-2xl sm:rounded-bl-2xl ${companyCommonStyles}`}
            >
              <p className="font-bold text-[20px]">Files</p>
            </div>
            <div className={companyCommonStyles}>
              <p className="font-bold text-[20px]">Time</p>
            </div>

            <div
              className={`sm:rounded-tr-2xl rounded-br-2xl ${companyCommonStyles}`}
            >
              <p className="font-bold text-[20px]">Download</p>
            </div>
          </div>
          {dataLength == 0 ? (
            <>
              <h1 className="text-3xl text-white  py-5">Loading... </h1>
              <Loader />
            </>
          ) : (
            [...serverFiles[0]].map((file, i) => (
              <div className="grid sm:grid-cols-3 w-full" key={i}>
                <>
                  <div
                    className={`rounded-tl-2xl sm:rounded-bl-2xl ${companyCommonStyles}`}
                  >
                    <p className="">{file}</p>
                  </div>

                  <div className={companyCommonStyles}>
                    {getTimestamp(file.toString())}
                  </div>
                  <div
                    className={`sm:rounded-tr-2xl rounded-br-2xl ${companyCommonStyles}`}
                  >
                    <li className="bg-[#2952e3] py-2 px-7 mx-4 rounded-lg cursor-pointer hover:bg-[#2546bd]">
                      <button
                        onClick={() => {
                          fetch(
                            `http://localhost:8080/download?filename=${file}`,
                          ).then((res) => {
                            res.blob().then((blob) => {
                              let url = window.URL.createObjectURL(blob)
                              let a = document.createElement('a')
                              a.style.display = 'none'
                              a.href = url
                              a.download = file
                              document.body.appendChild(a)
                              a.click()
                              window.URL.revokeObjectURL(url)
                            })
                          })
                        }}
                      >
                        {} Download
                      </button>
                    </li>
                  </div>
                </>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}

export default Archives
