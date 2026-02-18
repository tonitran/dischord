const COLORS = ['#5865f2', '#eb459e', '#57f287', '#fee75c', '#ed4245', '#3ba55c', '#faa61a']

interface Props {
  name: string
  size?: 'sm' | 'md' | 'lg'
}

const SIZES = {
  sm: 'w-8 h-8 text-xs',
  md: 'w-10 h-10 text-sm',
  lg: 'w-12 h-12 text-base',
}

export default function Avatar({ name, size = 'md' }: Props) {
  const color = COLORS[name.charCodeAt(0) % COLORS.length]
  return (
    <div
      className={`${SIZES[size]} rounded-full flex items-center justify-center text-white font-bold flex-shrink-0 select-none`}
      style={{ backgroundColor: color }}
    >
      {name[0].toUpperCase()}
    </div>
  )
}
