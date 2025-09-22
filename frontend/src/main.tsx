import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <h1 className="text-4xl text-pink-500">BeaconOfKnowledge!</h1>
    <App />
  </StrictMode>,
  
)
