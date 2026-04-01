import { useState } from 'react'
import { RocketIcon } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { PublishModal } from './PublishModal'

export const PublishButton = () => {
  const [showModal, setShowModal] = useState(false)

  return (
    <>
      <Button 
        variant="outline" 
        onClick={() => setShowModal(true)}
        className="gap-2"
      >
        <RocketIcon size={16} />
        Publish
      </Button>

      <PublishModal 
        isOpen={showModal} 
        onClose={() => setShowModal(false)} 
      />
    </>
  )
} 