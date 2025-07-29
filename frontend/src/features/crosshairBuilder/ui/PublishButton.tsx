import { useCrosshairStore } from '@/entities/crosshair/model/store'
import { RocketIcon } from 'lucide-react'
import { Button } from '@/shared/ui/button'

export const PublishButton = () => {
  const publish = useCrosshairStore(s => s.publish)
  return (
    <Button className="py-0 pe-0" variant="outline" onClick={publish}>
      <RocketIcon className="opacity-60" size={16} aria-hidden="true" />
      Publish
      <span className="text-muted-foreground before:bg-input relative ms-1 inline-flex h-full items-center justify-center rounded-full px-3 text-xs font-medium before:absolute before:inset-0 before:left-0 before:w-px">
        0
      </span>
    </Button>
  )
} 