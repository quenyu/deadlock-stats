import { CircleCheckIcon, X } from 'lucide-react'

interface Alert {
  id: string
  message: string
  type: 'svg' | 'config' | 'preset'
}

interface AlertsProps {
  alerts: Alert[]
  removeAlert: (id: string) => void
}

export const Alerts: React.FC<AlertsProps> = ({ alerts, removeAlert }) => (
  <div className="fixed bottom-4 right-4 z-50 space-y-2">
    {alerts.map((alert) => (
      <div
        key={alert.id}
        className="animate-in slide-in-from-right-2 duration-300 border border-border rounded-md px-4 py-3 bg-background shadow-lg"
      >
        <div className="flex items-center justify-between">
          <p className="text-sm flex items-center">
            <CircleCheckIcon
              className="me-3 -mt-0.5 inline-flex text-emerald-500"
              size={16}
              aria-hidden="true"
            />
            {alert.message}
          </p>
          <button
            onClick={() => removeAlert(alert.id)}
            className="cursor-pointer ml-2 text-muted-foreground hover:text-foreground transition-colors"
          >
            <X size={14} />
          </button>
        </div>
      </div>
    ))}
  </div>
) 