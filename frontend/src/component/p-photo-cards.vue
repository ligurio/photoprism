<template>
    <v-container grid-list-xs fluid class="pa-2 p-photos p-photo-cards">
        <v-card v-if="photos.length === 0" class="p-photos-empty secondary-light lighten-1 ma-1" flat>
            <v-card-title primary-title>
                <div>
                    <h3 class="title mb-3">
                        <translate>No photos matched your search</translate>
                    </h3>
                    <div>
                        {{$gettext("Try using other terms and search options such as category, country and camera.")}}
                        <span v-show="$config.feature('review')">
                            {{$gettext("Non-photographic and low-quality images require a review before they appear in search results.")}}
                        </span>
                    </div>
                </div>
            </v-card-title>
        </v-card>
        <v-layout row wrap class="p-results">
            <v-flex
                    v-for="(photo, index) in photos"
                    :key="index"
                    class="p-photo"
                    xs12 sm6 md4 lg3 d-flex
                    v-bind:class="{ 'is-selected': $clipboard.has(photo) }"
            >
                <v-hover>
                    <v-card tile slot-scope="{ hover }"
                            @contextmenu="onContextMenu($event, index)"
                            :dark="$clipboard.has(photo)"
                            :class="$clipboard.has(photo) ? 'elevation-10 ma-0 accent darken-1 white--text' : 'elevation-0 ma-1 accent lighten-3'">
                        <v-img :src="photo.thumbnailUrl('tile_500')"
                               aspect-ratio="1"
                               v-bind:class="{ selected: $clipboard.has(photo) }"
                               class="accent lighten-2 clickable"
                               @mousedown="onMouseDown($event, index)"
                               @click.stop.prevent="onClick($event, index)"
                        >
                            <v-layout
                                    slot="placeholder"
                                    fill-height
                                    align-center
                                    justify-center
                                    ma-0

                            >
                                <v-progress-circular indeterminate color="accent lighten-5"></v-progress-circular>
                            </v-layout>

                            <v-layout
                                    fill-height
                                    align-center
                                    justify-center
                                    ma-0
                                    class="p-photo-live"
                                    style="overflow: hidden;"
                                    v-if="photo.Type === 'live'"
                                    v-show="hover"
                            >
                                <video width="500" height="500" autoplay loop muted playsinline>
                                    <source :src="photo.videoUrl()" type="video/mp4">
                                </video>
                            </v-layout>

                            <v-btn v-if="hidePrivate && photo.Private" :ripple="false"
                                   icon flat large absolute
                                   class="p-photo-private opacity-75">
                                <v-icon color="white">lock</v-icon>
                            </v-btn>

                            <v-btn v-if="hover || selection.length && $clipboard.has(photo)" :ripple="false"
                                   icon flat large absolute
                                   :class="selection.length && $clipboard.has(photo) ? 'p-photo-select' : 'p-photo-select opacity-50'"
                                   @click.stop.prevent="onSelect($event, index)">
                                <v-icon v-if="selection.length && $clipboard.has(photo)" color="white"
                                        class="t-select t-on">check_circle
                                </v-icon>
                                <v-icon v-else color="accent lighten-3" class="t-select t-off">radio_button_off</v-icon>
                            </v-btn>

                            <v-btn icon flat large absolute :ripple="false"
                                   :class="photo.Favorite ? 'p-photo-like opacity-75' : 'p-photo-like opacity-50'"
                                   @click.stop.prevent="photo.toggleLike()">
                                <v-icon v-if="photo.Favorite" color="white" class="t-like t-on">favorite</v-icon>
                                <v-icon v-else color="accent lighten-3" class="t-like t-off">favorite_border</v-icon>
                            </v-btn>

                            <template v-if="photo.isPlayable()">
                                <v-btn v-if="photo.Type === 'live'" :ripple="false"
                                       icon flat large absolute class="p-photo-live opacity-75"
                                       @click.stop.prevent="openPhoto(index, true)" title="Live Photo">
                                    <v-icon color="white" class="action-play">adjust</v-icon>
                                </v-btn>
                                <v-btn v-else color="white" :ripple="false"
                                       outline large fab absolute class="p-photo-play opacity-75" :depressed="false"
                                       @click.stop.prevent="openPhoto(index, true)" title="Play">
                                    <v-icon color="white" class="action-play">play_arrow</v-icon>
                                </v-btn>
                            </template>
                            <v-btn v-else-if="photo.Type === 'image' && photo.Files.length > 1" :ripple="false"
                                   icon flat large absolute class="p-photo-merged opacity-75"
                                   @click.stop.prevent="openPhoto(index, true)">
                                <v-icon color="white" class="action-burst">burst_mode</v-icon>
                            </v-btn>
                            <v-btn v-else-if="photo.Type === 'raw'" :ripple="false"
                                   icon flat large absolute class="p-photo-merged opacity-75"
                                   @click.stop.prevent="openPhoto(index, true)" title="RAW">
                                <v-icon color="white" class="action-burst">photo_camera</v-icon>
                            </v-btn>
                        </v-img>

                        <v-card-title primary-title class="pa-3 p-photo-desc" style="user-select: none;">
                            <div>
                                <h3 class="body-2 mb-2" :title="photo.Title">
                                    <button @click.exact="editPhoto(index)">
                                        {{ photo.Title | truncate(80) }}
                                    </button>
                                </h3>
                                <div class="caption mb-2" v-if="photo.Description" title="Description">
                                    {{ photo.Description }}
                                </div>
                                <div class="caption">
                                    <button @click.exact="editPhoto(index)">
                                        <v-icon size="14" title="Taken">date_range</v-icon>
                                        {{ photo.getDateString() }}
                                    </button>
                                    <template v-if="!photo.Description">
                                        <br/>
                                        <button v-if="photo.Type === 'video'" @click.exact="openPhoto(index, true)"
                                                title="Video">
                                            <v-icon size="14">movie</v-icon>
                                            {{ photo.getVideoInfo() }}
                                        </button>
                                        <button v-else @click.exact="editPhoto(index)" title="Camera">
                                            <v-icon size="14">photo_camera</v-icon>
                                            {{ photo.getPhotoInfo() }}
                                        </button>
                                    </template>
                                    <template v-if="filter.order === 'name' && $config.feature('download')">
                                        <br/>
                                        <button @click.exact="downloadFile(index)"
                                                title="Name">
                                            <v-icon size="14">insert_drive_file</v-icon>
                                            {{ photo.baseName() }}
                                        </button>
                                    </template>
                                    <template v-if="showLocation && photo.LocationID">
                                        <br/>
                                        <button @click.exact="openLocation(index)" title="Location">
                                            <v-icon size="14">location_on</v-icon>
                                            {{ photo.getLocation() }}
                                        </button>
                                    </template>
                                </div>
                            </div>
                        </v-card-title>
                    </v-card>
                </v-hover>
            </v-flex>
        </v-layout>
    </v-container>
</template>
<script>
    export default {
        name: 'p-photo-cards',
        props: {
            photos: Array,
            selection: Array,
            openPhoto: Function,
            editPhoto: Function,
            openLocation: Function,
            album: Object,
            filter: Object,
        },
        data() {
            return {
                showLocation: this.$config.settings().features.places,
                hidePrivate: this.$config.settings().features.private,
                debug: this.$config.get('debug'),
                mouseDown: {
                    index: -1,
                    timeStamp: -1,
                },
            };
        },
        methods: {
            downloadFile(index) {
                const photo = this.photos[index];
                const link = document.createElement('a')
                link.href = "/api/v1/download/" + photo.Hash;
                link.download = photo.FileName;
                link.click()
            },
            onSelect(ev, index) {
                if (ev.shiftKey) {
                    this.selectRange(index);
                } else {
                    this.$clipboard.toggle(this.photos[index]);
                }
            },
            onMouseDown(ev, index) {
                this.mouseDown.index = index;
                this.mouseDown.timeStamp = ev.timeStamp;
            },
            onClick(ev, index) {
                let longClick = (this.mouseDown.index === index && ev.timeStamp - this.mouseDown.timeStamp > 400);

                if (longClick || this.selection.length > 0) {
                    if (longClick || ev.shiftKey) {
                        this.selectRange(index);
                    } else {
                        this.$clipboard.toggle(this.photos[index]);
                    }
                } else {
                    this.openPhoto(index, false);
                }
            },
            onContextMenu(ev, index) {
                if (this.$isMobile) {
                    ev.preventDefault();
                    ev.stopPropagation();
                    this.selectRange(index);
                }
            },
            selectRange(index) {
                this.$clipboard.addRange(index, this.photos);
            },
        }
    };
</script>
