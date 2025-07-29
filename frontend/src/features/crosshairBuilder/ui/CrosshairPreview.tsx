import React, { RefObject } from 'react'
import { CrosshairSettings } from '@/entities/crosshair/types/types'

interface CrosshairPreviewProps extends CrosshairSettings {
  isInteractive: boolean
  svgRef?: RefObject<SVGSVGElement>
  interactiveSvgRef?: RefObject<SVGSVGElement>
}

export const CrosshairPreview: React.FC<CrosshairPreviewProps> = ({
  color,
  thickness,
  length,
  gap,
  dot,
  opacity,
  pipOpacity,
  dotOutlineOpacity,
  isInteractive,
  svgRef,
  interactiveSvgRef,
}) => (
  <>
    {dot && <circle cx={60} cy={60} r={4} fill="#000" opacity={dotOutlineOpacity} />}
    {dot && <circle cx={60} cy={60} r={2} fill={color} opacity={opacity} />}
    <rect x={60 - thickness / 2} y={60 - gap / 2 - length} width={thickness} height={length} fill={color} opacity={pipOpacity} />
    <rect x={60 - thickness / 2} y={60 + gap / 2} width={thickness} height={length} fill={color} opacity={pipOpacity} />
    <rect x={60 - gap / 2 - length} y={60 - thickness / 2} width={length} height={thickness} fill={color} opacity={pipOpacity} />
    <rect x={60 + gap / 2} y={60 - thickness / 2} width={length} height={thickness} fill={color} opacity={pipOpacity} />
  </>
) 