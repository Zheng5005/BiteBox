import { useAuth } from "../context/AuthContext";

interface AvatarProps {
  size: Number
}

const Avatar: React.FC<AvatarProps> = ({size}) => {
  //TODO: MAKE THIS MORE DYNAMIC, IN ORDER TO USE IR IN DIFFERENT PAGES AND DIFFERENT CASES
  const { user } = useAuth();
  const stockImage = `https://avatar.iran.liara.run/public/${Math.floor(Math.random() * 100)}`
  let photo = user ? (user.url_photo != "" ? user.url_photo : stockImage) : stockImage

  return <>
    <img
      src={photo}
      alt="User"
      className={`avatar size-${size}`}
    />
  </>

}

export default Avatar
