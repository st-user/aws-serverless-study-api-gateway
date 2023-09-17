import { useState } from 'react'
import { Form } from './Form'


function App() {

  const [started, setStarted] = useState<boolean>(false)
  const [userId, setUserId] = useState<string>('')

  const isUserIdInput = userId.length > 0

  return (
    <>
      <h1>Sample</h1>
      <>
        <label>userId:</label>
        <input type="text" onChange={(e) => setUserId(e.target.value)}/>
      </>
      <>
        <button onClick={() => setStarted(true)} disabled={!isUserIdInput}>Run</button>
      </>
      {started && <Form userId={userId}/>}
    </>
  )
}

export default App
