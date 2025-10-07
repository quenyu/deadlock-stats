import React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Label } from '@/shared/ui/label'

interface PageSizeSelectorProps {
  pageSize: number
  onPageSizeChange: (pageSize: number) => void
  options?: number[]
}

export const PageSizeSelector: React.FC<PageSizeSelectorProps> = ({
  pageSize,
  onPageSizeChange,
  options = [10, 20, 50, 100]
}) => {
  return (
    <div className="flex items-center space-x-2">
      <Label htmlFor="page-size" className="text-sm font-medium">
        Show:
      </Label>
      <Select value={pageSize.toString()} onValueChange={(value) => onPageSizeChange(Number(value))}>
        <SelectTrigger className="w-20">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {options.map((size) => (
            <SelectItem key={size} value={size.toString()}>
              {size}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <span className="text-sm text-muted-foreground">per page</span>
    </div>
  )
} 