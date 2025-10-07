import { useState, useRef } from 'react'
import { useCrosshairStore } from '@/entities/crosshair/model/store'
import { PresetButtons } from './PresetButtons'
import { CrosshairPreview } from './CrosshairPreview'
import { CrosshairSettings } from './CrosshairSettings'
import { PublishButton } from './PublishButton'
import { Alerts } from './Alerts'
import previewImg from '@/shared/assets/images/preview_1.png'

interface Alert {
  id: string
  message: string
  type: 'svg' | 'config' | 'preset'
}

export const CrosshairBuilder = () => {
  const settings = useCrosshairStore(s => s.settings)
  const [isInteractive, setIsInteractive] = useState(false)
  const [alerts, setAlerts] = useState<Alert[]>([])
  const svgRef = useRef<SVGSVGElement>(null)
  const interactiveSvgRef = useRef<SVGSVGElement>(null)
  const posRef = useRef({ x: 60, y: 60 })
  const animating = useRef(false)

  const addAlert = (message: string, type: 'svg' | 'config' | 'preset') => {
    const id = Date.now().toString()
    setAlerts(prev => [...prev, { id, message, type }])
    setTimeout(() => {
      setAlerts(prev => prev.filter(alert => alert.id !== id))
    }, 3000)
  }

  const removeAlert = (id: string) => {
    setAlerts(prev => prev.filter(alert => alert.id !== id))
  }

  const updatePosition = () => {
    animating.current = false
    const { x, y } = posRef.current
    if (interactiveSvgRef.current) {
      interactiveSvgRef.current.style.transform = `translate(${x - 60}px, ${y - 60}px)`
    }
  }

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!isInteractive) return
    posRef.current = { x: e.nativeEvent.offsetX, y: e.nativeEvent.offsetY }

    if (!animating.current) {
      animating.current = true
      requestAnimationFrame(updatePosition)
    }
  }

  const handleMouseEnter = () => {
    setIsInteractive(true)
  }

  const handleMouseLeave = () => {
    setIsInteractive(false)
    posRef.current = { x: 60, y: 60 }
    if (interactiveSvgRef.current) {
      interactiveSvgRef.current.style.transform = 'translate(0px, 0px)'
    }
  }

  return (
    <div className="min-h-screen bg-background p-4">
      <div className="flex flex-col items-center justify-center min-h-screen">
        <Alerts alerts={alerts} removeAlert={removeAlert} />

        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent mb-2">
            Crosshair Builder
          </h1>
          <p className="text-slate-300 text-lg mb-6">Customize your Deadlock crosshair</p>
          
          <PresetButtons />
        </div>

        <div className="relative flex flex-col items-center justify-center w-full max-w-6xl mb-8">
          <div 
            className={`relative w-full max-w-6xl aspect-[16/9] flex items-center justify-center rounded-2xl overflow-hidden shadow-2xl border border-border bg-background/20 backdrop-blur-sm transition-all duration-300 hover:shadow-lg ${isInteractive ? 'cursor-none' : 'cursor-crosshair'}`}
            onMouseMove={handleMouseMove}
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
          >
            <img src={previewImg} alt="Deadlock preview" className="w-full h-full object-cover" />
            <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent"></div>
            
            {isInteractive ? (
              <svg
                ref={interactiveSvgRef}
                width="120"
                height="120"
                viewBox="0 0 120 120"
                className="absolute pointer-events-none drop-shadow-[0_0_12px_rgba(0,0,0,0.8)]"
                style={{ 
                  zIndex: 2,
                  left: 0,
                  top: 0
                }}
              >
                <CrosshairPreview {...settings} isInteractive={isInteractive} />
              </svg>
            ) : (
              <svg
                ref={svgRef}
                width="120"
                height="120"
                viewBox="0 0 120 120"
                className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 pointer-events-none drop-shadow-[0_0_12px_rgba(0,0,0,0.8)]"
                style={{ zIndex: 1 }}
              >
                <CrosshairPreview {...settings} isInteractive={false} />
              </svg>
            )}
          </div>
        </div>

        <CrosshairSettings />

        <div className="flex gap-4 mt-6">
          <PublishButton />
        </div>
      </div>
    </div>
  )
} 