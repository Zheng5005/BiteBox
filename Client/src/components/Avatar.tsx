import { useAuth } from "../context/AuthContext";

interface AvatarProps {
  size: number
}

const Avatar: React.FC<AvatarProps> = ({size}) => {
  const { user } = useAuth();
  const stockImage = `https://avatar.iran.liara.run/public/${Math.floor(Math.random() * 100)}`
  const photo = user ? (user.url_photo != "" ? user.url_photo : stockImage) : stockImage
  const px = size * 4;

  return (
    <img
      src={photo}
      alt="User"
      className="rounded-full object-cover"
      style={{ width: `${px}px`, height: `${px}px` }}
    />
  )
}

export default Avatar
