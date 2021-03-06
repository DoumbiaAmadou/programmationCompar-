 # -*- coding: utf-8 -*-
from case import Case
from ant import Ant
import pprint

class GameMap():

    def __init__(self, nb_ants, initial_energy, initial_acid,ia):
        self.map_dict = {}
        self.ants     = []

        for ant_id in range(0, nb_ants):
            ant = Ant(ant_id, initial_energy, initial_acid, True)
            self.ants.append(ant)

        self.ia=ia

    def __str__(self):
        text_repr = ""

        for coords, case in self.map_dict.iteritems():
            text_repr += "%s \t: %s\n" %(coords, case)

        return text_repr

    def get_ants(self):
        return self.ants

    def get_ant_with_id(self, ant_id):
        search_result = filter(lambda ant: ant.ant_id == ant_id, self.ants)
        if len(search_result) > 0:
            return search_result[0]
        else:
            return None

    def update(self, ants_stats, verbose=False):
        for stat in ants_stats:
            cases = map(lambda case_dict: Case(case_dict["content"]["kind"], int(case_dict["x"]), int(case_dict["y"])), stat[1])
            pprint.pprint(stat[1])
            for case in cases:
                self.map_dict[(case.x, case.y)] = case

            ant_state   = stat[0]
            ant_id      = ant_state["id"]
            current_ant = self.get_ant_with_id(ant_id)

            if current_ant is not None:
                current_ant.update_state(ant_state)
                if verbose:
                    print(current_ant)

        self.ia.play_turn()



        # if verbose:
            # print(self)

    def get_case_at(self, x, y):
        return self.map_dict.get((x, y), None)

    def get_case_surroundings(self, x, y):
        surroundings_coordinates = [(x-1, y-1),
                                    (x-1, y),
                                    (x-1, y+1),
                                    (x 	, y-1),
                                    (x 	, y+1),
                                    (x+1, y-1),
                                    (x+1, y),
                                    (x+1, y+1)]

        surrounding_cases = filter(lambda maybe_case: maybe_case is not None, map(lambda (x,y): self.get_case_at(x, y), surroundings_coordinates))

        return surrounding_cases

# save position of ennemies
    def get_position_enemies(self):
        enemies = []
        for case in map_dict:
            if map_dict[case] == "enemy":
                enemies.append(case)
        return enemies

# save position of food
    def get_position_resource(self):
        resources = []
        for case in map_dict:
            if map_dict[case] == "resource":
                resources.append(case)
        return resources

