import { Input } from '@/shared/ui/input'
import { Switch } from '@/shared/ui/switch'
import { Button } from '@/shared/ui/button'
import { Slider } from '@/shared/ui/slider'
import { useCrosshairStore } from '@/entities/crosshair/model/store'

export const CrosshairSettings = () => {
  const settings = useCrosshairStore(s => s.settings)
  const setSettings = useCrosshairStore(s => s.setSettings)

  const handleCopySVG = () => {
    // TODO: copy svg
    alert('Copy SVG (stub)')
  }

  const handleCopyConfig = () => {
    // TODO: copy config
    alert('Copy Config (stub)')
  }

  return (
    <div className="w-full max-w-4xl bg-card rounded-2xl shadow-2xl border border-border p-8">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="space-y-6">
          <div className="border-b border-border pb-2">
            <h3 className="text-xl font-semibold mb-4">Basic Settings</h3>
          </div>
          
          <div className="space-y-6">
            <div>
              <label className="block mb-3 font-medium">Color</label>
              <Input 
                type="color" 
                value={settings.color} 
                onChange={(e) => setSettings({ color: e.target.value })}
                className="h-12 w-full rounded-lg"
              />
            </div>
            
            <div>
              <label className="block mb-3 font-medium">Thickness</label>
              <div className="flex items-center gap-4">
                <Slider 
                  min={0} 
                  max={20} 
                  value={[settings.thickness]} 
                  onValueChange={value => setSettings({ thickness: value[0] })} 
                  className="flex-1" 
                />
                <span className="font-medium min-w-[3rem] text-center">{settings.thickness}px</span>
              </div>
            </div>
            
            <div>
              <label className="block mb-3 font-medium">Length</label>
              <div className="flex items-center gap-4">
                <Slider 
                  min={0} 
                  max={100} 
                  value={[settings.length]} 
                  onValueChange={value => setSettings({ length: value[0] })} 
                  className="flex-1" 
                />
                <span className="font-medium min-w-[3rem] text-center">{settings.length}px</span>
              </div>
            </div>
            
            <div>
              <label className="block mb-3 font-medium">Gap</label>
              <div className="flex items-center gap-4">
                <Slider 
                  min={-50} 
                  max={50} 
                  value={[settings.gap]} 
                  onValueChange={value => setSettings({ gap: value[0] })} 
                  className="flex-1" 
                />
                <span className="font-medium min-w-[3rem] text-center">{settings.gap}px</span>
              </div>
            </div>
            
            <div>
              <label className="block mb-3 font-medium">Pip Opacity</label>
              <div className="flex items-center gap-4">
                <Slider 
                  min={0} 
                  max={1} 
                  step={0.01} 
                  value={[settings.pipOpacity]} 
                  onValueChange={value => setSettings({ pipOpacity: value[0] })} 
                  className="flex-1" 
                />
                <span className="font-medium min-w-[3rem] text-center">{settings.pipOpacity.toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <div className="border-b border-border pb-2">
            <h3 className="text-xl font-semibold mb-4">Advanced Settings</h3>
          </div>
          
          <div className="space-y-6">
            <div className="p-4 rounded-lg bg-muted/50 border border-border">
              <div className="flex items-center gap-3 mb-3">
                <Switch 
                  checked={settings.dot} 
                  onCheckedChange={dot => setSettings({ dot })} 
                />
                <span className="font-medium">Center Dot</span>
              </div>
              {settings.dot && (
                <div className="mt-4 space-y-4">
                  <div>
                    <label className="block mb-3 font-medium">Dot Opacity</label>
                    <div className="flex items-center gap-4">
                      <Slider 
                        min={0} 
                        max={1} 
                        step={0.01} 
                        value={[settings.opacity]} 
                        onValueChange={value => setSettings({ opacity: value[0] })} 
                        className="flex-1" 
                      />
                      <span className="font-medium min-w-[3rem] text-center">{settings.opacity.toFixed(2)}</span>
                    </div>
                  </div>
                  <div>
                    <label className="block mb-3 font-medium">Dot Outline Opacity</label>
                    <div className="flex items-center gap-4">
                      <Slider 
                        min={0} 
                        max={1} 
                        step={0.01} 
                        value={[settings.dotOutlineOpacity]} 
                        onValueChange={value => setSettings({ dotOutlineOpacity: value[0] })} 
                        className="flex-1" 
                      />
                      <span className="font-medium min-w-[3rem] text-center">{settings.dotOutlineOpacity.toFixed(2)}</span>
                    </div>
                  </div>
                </div>
              )}
            </div>

            <div className="p-4 rounded-lg bg-muted/50 border border-border">
              <div className="space-y-4">
                <div className="flex items-center gap-3">
                  <Switch 
                    checked={settings.pipBorder} 
                    onCheckedChange={pipBorder => setSettings({ pipBorder })} 
                  />
                  <span className="font-medium">Show Pip Border</span>
                </div>
                
                <div className="flex items-center gap-3">
                  <Switch 
                    checked={settings.pipGapStatic} 
                    onCheckedChange={pipGapStatic => setSettings({ pipGapStatic })} 
                  />
                  <span className="font-medium">Static Pip Gap</span>
                </div>
                
                <div>
                  <label className="block mb-3 font-medium">Hit Marker Duration</label>
                  <div className="flex items-center gap-4">
                    <Slider 
                      min={0} 
                      max={1} 
                      step={0.01} 
                      value={[settings.hitMarkerDuration]} 
                      onValueChange={value => setSettings({ hitMarkerDuration: value[0] })} 
                      className="flex-1" 
                    />
                    <span className="font-medium min-w-[3rem] text-center">{settings.hitMarkerDuration}s</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <Button 
              onClick={handleCopySVG}
              className="flex-1 cursor-pointer"
            >
              Copy SVG
            </Button>
            <Button 
              variant="secondary" 
              onClick={handleCopyConfig}
              className="flex-1 cursor-pointer"
            >
              Copy Config
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
} 