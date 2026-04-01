import { Button } from '@/shared/ui/button'
import { PRESETS } from '@/entities/crosshair/lib/presets'
import { useCrosshairStore } from '@/entities/crosshair/model/store'

export const PresetButtons = () => {
  const loadPreset = useCrosshairStore(s => s.loadPreset)
  return (
    <div className="flex flex-wrap justify-center gap-3 mb-6">
      {Object.keys(PRESETS).map((presetName) => (
        <Button
          key={presetName}
          variant="outline"
          size="sm"
          onClick={() => loadPreset(presetName)}
          className="capitalize"
        >
          {presetName}
        </Button>
      ))}
    </div>
  )
} 