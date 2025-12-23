import * as React from "react"
import { X } from "lucide-react"
import { Badge } from "./badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./select"

interface MultiSelectProps {
  options: string[]
  value: string[]
  onChange: (value: string[]) => void
  placeholder?: string
}

export function MultiSelect({ options, value, onChange, placeholder = "Select items..." }: MultiSelectProps) {
  const [open, setOpen] = React.useState(false)

  const handleSelect = (selectedValue: string) => {
    if (!value.includes(selectedValue)) {
      onChange([...value, selectedValue])
    }
    setOpen(false)
  }

  const handleRemove = (itemToRemove: string) => {
    onChange(value.filter((item) => item !== itemToRemove))
  }

  const availableOptions = options.filter((option) => !value.includes(option))

  return (
    <div className="space-y-2">
      <Select open={open} onOpenChange={setOpen} onValueChange={handleSelect}>
        <SelectTrigger>
          <SelectValue placeholder={placeholder} />
        </SelectTrigger>
        <SelectContent>
          {availableOptions.length === 0 ? (
            <div className="px-2 py-1.5 text-sm text-gray-500">No options available</div>
          ) : (
            availableOptions.map((option) => (
              <SelectItem key={option} value={option}>
                {option}
              </SelectItem>
            ))
          )}
        </SelectContent>
      </Select>

      {value.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {value.map((item) => (
            <Badge key={item} variant="secondary" className="flex items-center gap-1">
              {item}
              <button
                type="button"
                onClick={() => handleRemove(item)}
                className="ml-1 rounded-full hover:bg-gray-700"
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  )
}
