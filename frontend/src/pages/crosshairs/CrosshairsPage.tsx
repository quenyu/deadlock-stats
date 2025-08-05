import { useState } from 'react'
import { CrosshairGallery } from '@/features/crosshairGallery/ui/CrosshairGallery'
import { PublishModal } from '@/features/crosshairBuilder/ui/PublishModal'
import { Button } from '@/shared/ui/button'
import { PlusIcon } from 'lucide-react'

export const CrosshairsPage = () => {
  const [showPublishModal, setShowPublishModal] = useState(false)

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold mb-2">Crosshair Gallery</h1>
          <p className="text-muted-foreground">
            Explore and rate crosshairs from the Deadlock community
          </p>
        </div>
        <Button 
          onClick={() => setShowPublishModal(true)}
          className="gap-2"
        >
          <PlusIcon size={16} />
          Publish crosshair
        </Button>
      </div>

      <CrosshairGallery />

      <PublishModal 
        isOpen={showPublishModal} 
        onClose={() => setShowPublishModal(false)} 
      />
    </div>
  )
} 