import { useState } from 'react'
import { useCrosshairStore } from '@/entities/crosshair/model/store'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/shared/ui/alert-dialog'

interface PublishModalProps {
  isOpen: boolean
  onClose: () => void
}

export const PublishModal: React.FC<PublishModalProps> = ({ isOpen, onClose }) => {
  const publish = useCrosshairStore(s => s.publish)
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [isPublic, setIsPublic] = useState(true)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handlePublish = async () => {
    if (!title.trim()) {
      setError('Title is required')
      return
    }

    setLoading(true)
    setError('')

    try {
      await publish(title.trim(), description.trim(), isPublic)
      setTitle('')
      setDescription('')
      setIsPublic(true)
      onClose()
    } catch (err) {
      setError('Failed to publish crosshair')
    } finally {
      setLoading(false)
    }
  }

  // @ts-expect-error - Will be used for proper modal close handling
  const handleClose = () => {
    if (!loading) {
      setTitle('')
      setDescription('')
      setIsPublic(true)
      setError('')
      onClose()
    }
  }

  return (
    <AlertDialog open={isOpen} onOpenChange={() => { if (!loading) onClose(); }}>
      <AlertDialogContent className="max-w-md">
        <AlertDialogHeader>
          <AlertDialogTitle>Publish Crosshair</AlertDialogTitle>
          <AlertDialogDescription>
            Share your crosshair with the community. Fill in the title and description.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <div className="space-y-4 py-4">
          <div className="space-y-2">
            <label htmlFor="title" className="text-sm font-medium">
              Title *
            </label>
            <Input
              id="title"
              placeholder="My awesome crosshair"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={loading}
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="description" className="text-sm font-medium">
              Description
            </label>
            <Textarea
              id="description"
              placeholder="Describe the features of your crosshair..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={loading}
              rows={3}
            />
          </div>

          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="isPublic"
              checked={isPublic}
              onChange={(e) => setIsPublic(e.target.checked)}
              disabled={loading}
              className="rounded border-gray-300"
            />
            <label htmlFor="isPublic" className="text-sm">
              Make public
            </label>
          </div>

          {error && (
            <div className="text-sm text-red-500 bg-red-50 p-2 rounded">
              {error}
            </div>
          )}
        </div>

        <AlertDialogFooter>
          <AlertDialogCancel disabled={loading}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={handlePublish}
            disabled={loading || !title.trim()}
            className="bg-primary text-primary-foreground hover:bg-primary/90"
          >
            {loading ? 'Publishing...' : 'Publish'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
} 