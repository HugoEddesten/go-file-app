import { Link } from "react-router-dom"
import { Separator } from "../../components/ui/separator"

export const Nav = () => {

  return (
    <div className="w-full h-full flex flex-col">
      <div className="flex text-center gap-4 items-center p-2">
        <Link to={"/"} className="w-12">
          <img src="../../../public/SecureArchive-favicon.png"/>
        </Link>
        <Separator orientation="vertical"/>
        
        <Link to={"/"}>
          Home
        </Link>
        <Separator orientation="vertical"/>

        <Link to={"/profile"}>
          Profile
        </Link>
      </div>
      <Separator />
    </div>
  )
}