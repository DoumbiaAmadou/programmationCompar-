import scalaj.http._

object ConnectApi {
	private val baseUrl = "http://yann.regis-gianas.org/antroid/"
	private val version = "0"

	def makeHttpRequest(requestName: String): HttpRequest = {
		val request: HttpRequest = Http(baseUrl+version+"/"+ requestName +"?")
		return request
	}  
	// Send request by method POST
	def postRequest(requestName: String, params: Seq[(String, String)]) : String = {
		val request: HttpRequest = Http(baseUrl+version+"/"+ requestName +"?")
      	val response: HttpResponse[String] = request.postForm(params).asString
      	return response.body
	}

	// Send request by method GET (just for the request hasn't param)
	def getRequest(requestName: String) : String = {
		val request: HttpRequest = Http(baseUrl+version+"/"+ requestName +"?")
      	val response1: HttpResponse[String] = request.asString
      	return response1.body
     }


	// Connect a user
	def auth(login: String, password: String) : String ={
		var params = Seq("login" -> login, "password" -> password)
		return postRequest("auth",params)
	}

	// Create a new game. -- probleme: can't send request with params int
	//	def newGame(users: String, teaser: String, pace: String , nb_turn: Int, nb_ant_per_player: Int, nb_player: Int,
		//			 minimal_nb_player: Int, initial_energy: Int, initial_acid: Int) : String = {
		//	val request: HttpRequest = Http(baseUrl+version+"/create?")
		//	var postform = Seq("users"->users,"teaser"->teaser,"pace"->pace, "nb_turn"->nb_turn, 
		//						"nb_ant_per_player"->nb_ant_per_player, "nb_player"->nb_player, 
		//						"minimal_nb_player"->minimal_nb_player, "initial_energy"->initial_energy,
		//						"initial_acid"->initial_acid)
      	//	val response: HttpResponse[String] = request.params(postform).asString
      	//	return response.body
	//	}

	// Destroy a game.
	def destroyGame(idGame: String) : String = {
		val request: HttpRequest = Http(baseUrl+version+"/destroy?")

      	// val response: HttpResponse[String] = makeHttpRequest("destroy").param("id", idGame).asString
      	val response: HttpResponse[String] = Http(baseUrl+version+"/destroy?").param("id","1").asString
		return response.body
	}

	// Join a game.
	def joinAGame(idGame: String): String = {
		val request: HttpRequest = Http(baseUrl+version+"/join?")

      	// val response: HttpResponse[String] = makeHttpRequest("destroy").param("id", idGame).asString
      	val response: HttpResponse[String] = request.param("id",idGame).asString
		return response.body
	}

	// List all visible games.
	def listGame(): String = {
		return getRequest("games")
	}

	// Register a user.
	def registerUser(login: String, password: String): String = {
		var params = Seq("login" -> login, "password" -> password)
		return postRequest("register",params)
	}
	
	// Test connect
	def main(args: Array[String]) {
      	    println("Hello, scala request!")
	    // println(auth("nga","ngango"))
		// println(newGame("nga", "ngo", 10,30,1,1,1,100,100))
		println(listGame())
		// println(destroyGame("4986356145915714886550832111036009180"))
		// println(joinAGame("4986356145915714886550832111036009180"))

    }

}
